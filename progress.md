# Progress 5.26 (develop a language)


## Current
- added functions to compiler and VM
#### some data structures used
in compiler
``` go
type SymbalTable struct {
	outer *SymbalTable
	store map[string]*Symbol
	size  int
}

func NewEnclosedSymbalTable(outer *SymbalTable) *SymbalTable {
	st := NewSymbolTable()
	st.outer = outer
	return st
}

func (s *SymbalTable) Define(name string) *Symbol {
	symbol := Symbol{Name: name, Scope: GlobalScope, Index: s.size}
	if s.outer != nil {
		symbol.Scope = LocalScope
	}
	s.store[name] = &symbol
	s.size++
	return &symbol
}
```

in virtual machine
```go
type Frame struct {
	fn *CompiledFunction
	ip int
	bp int
}
type CompiledFunction struct {
	Instructions Instructions
	NumLocals    int
	NumParas     int
}
```
local variables are stored on stack
```go
3 : nil // top of stack
2 : 2
1 : 1 // base pointer
0 : add() // result would be placed here
```
#### result
```go
	fib := func(n) {
		if n <= 2 {
			n
		} else {
			fib(n-1) + fib(n-2)
		}
	}
	fib(35)  // 40s...very slow
```
## language feature
- hope 
  - a keyword to check correctness, particularly when we define a function
  - we give parameters and expected answer
  
```go
	fib := func (n) {
		if n <= 2 {
			n
		} else {
			fib(n-1) + fib(n-2)
		}
	} hope {
		1 -> 1
		2 -> 2
		3 -> 3
	}
```
```go
	add := func(x, y) {
		x + y
	} hope {
		1,2 -> 3
		-1,0 -> -1
	}
```
```go
	//hope: 1 -> 1
	//hope: 2 -> 2
	//hope: 3 -> 3
	fib := func (n) {
		//...
	}
```
```go
	//hope: 1,2 -> 3
	//hope: 2,4 -> 6
	//hope: 3,3 -> 6
	add := func (x, y) {
		//...
	}
```
