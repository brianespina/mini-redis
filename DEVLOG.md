# mini-redis devlog
A build and learn Redis clone in Go. The primary goal is to understand every layer  form socket up. this log captures  thinking-as-it-happened: what broke, what confused me,  what I decided and why. 

--

# how to use this log
One enty per session. 

- **what  i built** - one or two lines
- **what broke** - bugs, dead-ends, wrong models. 
- **what  I decided and why** - especially forks where I picked one option over another.
- **concept unlocked** - arguably the main goal of this project. to unlock concepts.
- **for future me** - whatever I want to say to myself on the next session. 

Keep it honest

--

## Phase 1 and 2 - TCP transport

**Built**  A TCP server on :6379 that accepts connections, spawns  a goroutine per client , reads saw bytes in a loop,  cleans up  disconnect.

**Decided:**
- One goroutine per connection (`go handleConn(conn)`) — idiomatic Go, and it's what frees the accept loop to take the next client immediately.
- Port 6379 just so real Redis clients connect with zero config. Arbitrary otherwise (it spells a name on a phone keypad, apparently).


## Phase 2 -  RESP parsing 

**Built** A Parser that parses RESP  wire format to `[]string`

**What broke / confused me**
- First instict is to **split the string on `\r\n`**. Its wrong, values are binary safe and can also contain `\r\n`. also you can't split the stream, it may still be arriving. 
- Fed `1\r\n` to `Atoi`, apparently ReadString('\n') keeps the terminator. 
- Forgot to consume trailing `\r\n` after payload, so desync happened. Miserable to debug until I understood the mechanism.

**Concept Unlocked** (a lot on this first session)
- **Streams** this is totally new to me, I was treating it as a finished string. so its data incoming, still arriving, unfinished, or finished.   
- **Reading Streams** so reading streams is like reading whatever is available at the moment. Every reaad moves the cursor forward, permanently. There are a bunch of types of Read function.
- **Desync** is a thing, its you code expecting one thing and the cursor is on a different place, so its just fucks up the data, and the code behavior.

## Phase 3 — The store + serializing (the write side)
**Built:** `SET`, `GET`, `ECHO` against an in-memory `map[string]string`. Write-side helpers to produce RESP: bulk string (`$%d\r\n%s\r\n`), null (`$-1\r\n`), simple string (`+OK\r\n`), error (`-ERR ...\r\n`).
**What broke / confused me:**
- GET double-wrote its reply: sent the value AND the null, because the found-branch had no `break`/`else`. The stray `$-1\r\n` desynced the *client* — every subsequent reply answered the wrong command. **Write-side desync**, the mirror of the parser bug. One missing `else` made three commands look broken.
- Nearly spliced raw client input into an error reply (`-ERR unknown command %s`). If the input contained `\r\n`, that breaks the client's framing → protocol injection, same family as SQL injection. Error/simple strings have no length prefix, so they're the *only* place embedded `\r\n` can break framing.
**Decided:**
- Missing key → `$-1\r\n` (null), **not** an error. A missing key is a valid answer of "nothing," distinct from "your request was wrong." Represent absence properly — every system has to (SQL NULL, JS undefined, Go's `ok`).
- Use `v, ok := store[key]` to distinguish "missing" from "stored empty string."
- Data lives in RAM, lost on restart — and that's the *point*. Redis is deliberately in-memory for speed. Durability comes later (AOF), optional.
- On a **write error**, don't reply — the connection's already broken. Just `return` and let cleanup happen. Replying down a dead pipe is futile.
**Concept unlocked:** The write side is the trivial inverse of the read side — just string formatting with length prefixes. The length prefix I *trusted* while reading is the length prefix I now *promise* while writing. Binary-safe both ways.
**For future-me:** Serializer helpers live in `resp.go`, write-side twins of the read helpers. That file is now the single place that knows what RESP bytes look like, in both directions.
---

## Phase 4 — Concurrency safety
**Built:** Moved `store` to package level (shared by all goroutines), then guarded it with a mutex. Reads use `RWMutex.RLock`, writes use `Lock`.
**What broke / confused me:**
- With `store` *inside* `handleConn`, each client got a private map (local var, own stack) → two windows couldn't see each other's data. Moved it to package level. Fixed the isolation... and introduced a far worse bug.
- Shared map, no lock → **data race**. Tricky part: it ran clean for *3 million* operations before I could make it fail. Concurrency bugs are *probabilistic* — "it worked" proves nothing. The bug was in the source the whole time.
- `go run -race` on the real server showed nothing at first, because a single connection = a single goroutine = no concurrency. Needed *multiple persistent connections at once* to force overlap. Isolated the detector with a tiny racetest.go (guaranteed race) to confirm the tool worked — it did.
- Once overlap was real: `fatal error: concurrent map writes`. That's the Go *runtime's* always-on guard (write-write only), separate from `-race` (opt-in, broader, catches read-write too). The runtime guard fires first and kills the process, which is why I rarely saw the `-race` WARNING — the harsher alarm wins the footrace.
**Decided:**
- `sync.Mutex` first (correct, simple), then upgraded to `sync.RWMutex` because the store is read-heavy — many GETs can `RLock` in parallel; only SETs need exclusive `Lock`. Applied by reasoning, not a visible benchmark.
- **Lock narrow.** Guard only the map access, never the `c.Write`. Copy the value out under the lock, release, *then* do network I/O. Holding a lock across slow I/O would serialize every client behind one connection.
- **Single exit for the lock.** Flattened GET to lock→copy→unlock→branch, so there's exactly one `Unlock` and no code path can skip it (that'd deadlock the whole server forever).
**Concept unlocked:**
- **The goroutine → thread → core tower.** Many cheap goroutines scheduled onto few OS threads (Go's scheduler) scheduled onto few cores (OS scheduler). "Parked" = suspended with a bookmark, 0% CPU. Blocking is cheap because parked goroutines don't hold a thread.
- **Node vs Go** are the same epoll machinery, opposite steering wheels: Node hands you events (callbacks), Go hides them in the scheduler so I can write plain blocking code. The scheduler is the price of that illusion.
- **Zero values.** `var mu sync.Mutex` needs no `=` — a zero-value mutex is usable. But `var m map[...]` is `nil` and unusable, so `store` needs `= map[...]{}`. Knowing which types are usable-at-zero *is* knowing the type.
- **Two frames for concurrency:** (1) guard the shared thing (mutex), or (2) give it one owner and pass messages (channels — "share memory by communicating"). Real Redis is single-threaded = the channel/owner model. I used a mutex because it's simple shared state; that's a legit, idiomatic choice, not the "primitive" one.
**For future-me:** The store is concurrency-safe *by construction* now — no amount of hammering can reproduce the crash, because the lock makes the collision impossible, not just improbable. That's the difference between "passed the test" and "cannot fail." Channels are coming for pub/sub, where message-passing fits the problem naturally.
---

## Recurring theme: the missing `return`
The disconnect error check (`if err != nil { print; return }`) lost its `return` roughly **four times** across refactors. Without it, a dead connection spins at 100% CPU forever. Lesson banked: *checking* an error is nothing; the decision *after* it is the handling. Made the disconnect test (`ctrl-C` the client, expect exactly one log line) a permanent part of the gauntlet so tooling catches what my eyes miss.
---

## The gauntlet (regression test, run after every change)
- `PING` → `PONG`
- `FOO` → `(error) ERR unknown command`
- `ECHO hello` → `"hello"`; `ECHO "a b c"` → `"a b c"` (space stays one arg)
- `SET k v` → `OK`; `GET k` → `"v"`; `GET missing` → `(nil)`
- `ctrl-C` the client → exactly ONE "client disconnected" line, no CPU spin
- `printf 'garbage\r\n' | nc localhost 6379` → server drops connection with a clear parse error in the log (the validation tripwire firing)
- Under load: 3 terminals hammering SET/GET, `go run -race .` → stays green
---

## Roadmap (next base camps)
- [ ] **Typed values** — `map[string]string` → values that can be string OR list OR hash. Forces Go **interfaces** and my first custom type. WRONGTYPE errors.
- [ ] **Expiry** — keys with TTLs. Lazy expiry (delete on read) first, then an active background sweep. Time as state.
- [ ] **AOF persistence** — append each write command's bytes to a file; on boot, replay the file *through my own parser*. (The `io.Reader` parser already does this — feed it a file instead of a socket.)
- [ ] **Transactions** — `MULTI`/`EXEC`, command queuing, atomicity.
- [ ] **Pub/sub** — fan-out to subscribers. Natural fit for **channels** — revisit the message-passing frame here.
---

## Reading list (compare my version to the real thing)
- The RESP spec (short) — confirm I'm handling the types right.
- antirez-era Redis C source — famously readable. The gap between my version and theirs is where the scale lessons live.
- "Go's race detector: the bugs it misses" — dynamic detection is probabilistic; a green `-race` run is not proof of safety.

