The server listens on `:6379`. Connect with any Redis client:
 
```sh
valkey-cli -p 6379
# or: redis-cli -p 6379
```
 
```
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> SET name Brian
OK
127.0.0.1:6379> GET name
"Brian"
127.0.0.1:6379> GET missing
(nil)
```
 
## Commands
 
| Command | Description |
|---------|-------------|
| `PING` | Replies `PONG`. |
| `ECHO <msg>` | Replies with `<msg>`. |
| `SET <key> <value>` | Stores a string value. |
| `GET <key>` | Returns the value, or `(nil)` if the key doesn't exist. |
| `COMMAND` | Returns an empty reply (enough to satisfy client handshakes). |
 
## What it does
 
- Accepts multiple concurrent clients, one goroutine per connection.
- Parses the RESP wire protocol (binary-safe, length-prefixed).
- Serializes RESP replies (simple strings, bulk strings, nulls, errors).
- Stores data in an in-memory map, shared safely across connections with a
  `sync.RWMutex` (reads can run in parallel; writes are exclusive).
- Data lives in RAM and is lost on restart — the same in-memory design Redis
  itself uses for speed. Persistence is planned (see roadmap).
## Tests
 
```sh
go test ./...          # run the suite
go test -v ./...       # verbose, shows each case
go test -race ./...    # run with the data-race detector
```
 
Table-driven unit tests cover the RESP parser (`resp_test.go`), including
malformed input, truncated streams, empty values, and binary-safe values that
contain `\r\n`.
 
## Project layout
 
```
main.go        entry point, TCP accept loop
handler.go     per-connection loop and command dispatch
resp.go        RESP protocol: reading and writing
resp_test.go   parser tests
```
 
## Roadmap
 
- [x] TCP server, one goroutine per connection
- [x] RESP parsing and serializing
- [x] `SET` / `GET` string store
- [x] Concurrency-safe store (`RWMutex`)
- [x] Unit tests for the parser
- [ ] Typed values — lists, hashes, sets (needs interfaces; `WRONGTYPE` errors)
- [ ] Key expiry / TTL
- [ ] AOF persistence (append write commands to disk, replay on boot)
- [ ] Transactions (`MULTI` / `EXEC`)
- [ ] Pub/sub
