# Ho
A programming language 
- named Ho 
- written in Go 
- inspired by Go and Pyret
- based on interpreter and compiler

Ho has features:
- inline testing
  

Ho is like
```Go
fib := func(n int) {
  n <= 2 ? n : fib(n-1) + fib(n-2)
} hope {
  1 -> 1
  2 -> 2
  3 -> 3
  10 -> 89
}
```
hope block is like:
```Go
func (...) {
  ...
} hope {
  input -> expected output
  ...
}
```


