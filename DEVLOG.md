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

---

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

---

## Log 6: Unit Testing 07/20/2026
**Built** Table-driven tests cases for `readInt` and `readBulkString`. 

**What broke / confused me**
- the `wantErr` boolean is a just a bit confusing to set. 
- the logic of one test is wrong, `wantErr` needed to be true not false. 

**Concept Unlocked**
- I understanding the structure of tests more and more. the mental model is starting to click. 
- writing test is like essentially you trying to dictate the function's behavior on how to handle inputs, good or bad, I now understand why some developers start with Tests before writing the function, I've heard this is called test-driven development.

**For future-me** do more reps still writing tests, the goal before moving forward with the project is to write tests for all the functions in the RESP parser. 

---

## Log 7: Unit Testing 07/21/2026
**Built** Created tests for all of the parser functions. 

**What broke / confused me**
- I'm still getting confused with the conditional statements, I said `if slices.Equal` then fail the test. it should be `if !slices.Equal` then fail the test. the logic is inverted in tests so I'm still getting tripped up with that.  
- got confused on the concept of testing finite `io.Reader`s like `strings.NewReader` and actual client connections without EOF. 
- Not sure if my test cases are good, I know some of them are overlapping on some functions, I'm using `readInt` in `readCommands` so there is some overlap in test cases, is that ok? or is that bad practice? I don't know. 

**Concept Unlocked**
- Stream gets clearer while writing the tests, when writing test I use `strings.NewReader` this simulates a stream, but the thing with this is that `strings.NewReader` returns an EOF. in a true connection to a client, the goroutine gets parked waiting, and that is the correct behaviour, to wait for more commands. now my next question is do you just keep waiting forever? apparently in production no, there is this thing called "slowloris", its essentially expiring the connection is how I understand it, closing the conection after some time of silence. 
- so apparently a little overlap is fine, but ideally minimal, you have to check if for instance `readCommands` is introducing a point of failure, something like that. 

**For future-me** I'm ending this Unit Testing Phase for now, since I finished the Goal of creating tests for all the functions in the RESP parser. I'll keep writing tests for my functions moving forward. Testing is still new to me but its not black box anymore. one conclusion I've made with this is that, Testing is a skill of its own that I need to have the muscle for if I want to be a serious programmer. 

---


## Log 8: Storing multiple types 07/23/2026
**Built** Implemented LPUSH command.

**What broke/ confused me**
- I'm kinda new at type assertions, I encountered it before I did not really understand it that much. 
- confused working with `any` type. 
- so in LPUSH you are reading updating and then writing, I wrote the locks only just encapsulating the actual read and writes. thought it would be easier to see and not miss unlocks. but apparently this intoduces errors when 2 clients are doing LPUSH. 

**Concept Unlocked** 
- ok type assertions, I understand them more now in the contect of having `any` type, you are just saying to the compiler, this `any` "right here I expect this type, is it?" , it returns two things, the first one, the value(or pointer?) as the type you asserted with, or nil if its not?, and the second a `bool` if the assertion was successfull? or possible? I need to look into this more. 
- So working with the locks I had a deeper understanding of `slices`, `arrays` because of the function `append()`. so slices are a structs that has a pointer, a lenth, and a capacity. I vaguely remember this when I was studying C, I had to implement append. when you initialize a `slice` or create a slice varaiable, or read a slice. Whatever the variable its being hel at also has a new copy of the length and capacity, the variable has its own bookeeping of those numbers but has the pointer to the same memmory, so if there are 2 clients holding their own separate bookeeping of the same memory writes at the same time, that is going to introduce problems. so the correct thing to do: whenever you are doing read-update-write on the shared store, you need to treat the whole opperation as one whole thing and lock it.  

**For future-me**: 
- look more into assertions, now that the store is type  
- `map[string]any` GET is broken. fix that next. 
- this implementation of storing multiple types, I can already see as very flimsy. I know that. I decided to do this implementation because this is exactly how I would implement it. I want to feel how this sucks, I wan't to feel how flimsy it is, after a couple more command implementation. I'm going to study the right way and refactor the commands. 

---

## Log 9: LRANGE 07/24/2026
**Built ** Implemented LRANGE command. 

**What broke/ confused me**
- the logic of the range itself, reading the redis documentation, I saw `LRANGE key 0 -1`, I see -1 and imidiately thought `oh the indexes loops!` and so I used modulo, I implemented it `start = start % length` `end = end % length` this gave me all sorts of unexpected results. including panicing the server. 
- of-by-one error in my clamp logic. working with arrays indexes and lengths are still so confusing to me. 

**Concept unlocked**
- with Go slices you can go over the len for instance an array with `len = 3` and you go `array[0:4]` <- this is valid as long as its not over the capacity. this triped me up debugging the of-by-one. 
- figured out that the indexes are not looping. there is no modulo at all. you just calculate the offset by adding the lenght if there are negative numbers, then clamp them. 

**For future-me:**
- do not move on from this command until you have a mental pictuse of the of-by-one, and it makes sense in your head. 
- you need more reps of working with arrays. (maybe leet code problems)
- there is still an of-by-one bug in the code, fix it
