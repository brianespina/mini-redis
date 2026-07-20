# mini-redis devlog
A build and learn Redis clone in Go. The primary goal is to understand every layer  form socket up. this log captures  thinking-as-it-happened: what broke, what confused me,  what I decided and why. 

---

# how to use this log
One enty per session. 

- **what  i built** - one or two lines
- **what broke** - bugs, dead-ends, wrong models. 
- **what  I decided and why** - especially forks where I picked one option over another.
- **concept unlocked** - arguably the main goal of this project. to unlock concepts.
- **for future me** - whatever I want to say to myself on the next session. 

Keep it honest

---

# Note
I began keeping this log while writing tests for the parser; everything before that lives in the commit history and my memory

## Log 5: Unit Testing 07/20/2026
**Built** my first test file for the `readInt` function of the parser. `resp_test.go`. 

**What broke / confused me**
- The whole mental model of writing test is totally different. the logic got jumbled in my head. its a totaly new skill that I'm excited to learn. 
- The difference of `t.Fatal` and `t.Error`. I just wrote the test and did not really understand what I'm writing or if its even usefull.
- How do you even test it without running server, a client connected and sending commands. 

**Concept Unlocked**
- I started to wrap my head around the pattern(or maybe structure is the better word) of tests. Its still a bit confusing at times but I just need reps
- `t.Fatal` totaly stops the the test and fails it. `t.Error` does not stop the test, it takes note of the error and keeps going. 
- so you can simulate a stream with `string.NewReader`. anything `io.Reader` interface can is stream like, be it a string or file or socket. anything pull-based, read-whatever's-availablet-now source. 

**For future-me** do more reps with writing tests, write table-based(whatever that is) test. before you move on to anything else in the project. write tests for the parser. 

## Log 6: Unit Testing 07/20/2026
**Built** Table-driven tests cases for `readInt` and `readBulkString`. 

**What broke / confused me**
- the `wantErr` boolean is a just a bit confusing to set. 
- the logic of one test is wrong, `wantErr` needed to be true not false. 

**Concept Unlocked**
- I understanding the structure of tests more and more. the mental model is starting to click. 
- writing test is like essentially you trying to dictate the function's behavior on how to handle inputs, good or bad, I now understand why some developers start with Tests before writing the function, I've heard this is called test-driven development.

**For future-me** do more reps still writing tests, the goal before moving forward with the project is to write tests for all the functions in the RESP parser. 

