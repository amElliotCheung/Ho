# Progress 6.1 (develop Ho language)
## Finished
#### changed ip to a local variable in vm.Run()
``` go
	ip := 0
	frame := vm.currentFrame()
	ins := frame.fn.Instructions
	
	for ip < len(ins) {
		op := Opcode(ins[ip])
		switch op {
			//......

		case OpCall:
			nextFrame := NewFrame(fn, ip+2, vm.stackIdx-numParas)
			//......
			vm.pushFrame(nextFrame)
			ip, frame = 0, vm.currentFrame()

		case OpReturnValue:
			//......
			f := vm.popFrame()
			vm.stack[f.bp-1] = result
			vm.stackIdx = f.bp
			ip, frame = f.ip, vm.currentFrame()
			ins = frame.fn.Instructions
		}
	}
```

#### added keyword "hope"
- add an instruction "OpHope"

- changed AST
```go
type FunctionLiteral struct {
	Parameters []*IdentifierLiteral
	Execute    *BlockExpression
	Hopes      *HopeBlock
}
type HopeBlock struct {
	HopeExpressions []HopeExpression
}
type HopeExpression struct {
	Parameters []Expression
	Expected   Expression
}
```
When we compile a function literal
```go
	for i, hopeExpr := range node.Hopes.HopeExpressions {
			c.emit(OpConstant, idx)
			for _, para := range hopeExpr.Parameters {
				c.Compile(para)
			}
			c.emit(OpCall, len(node.Parameters))
			c.Compile(hopeExpr.Expected)
			// 1 means the length of the expected answer
			// for now, function returns only one value
			c.emit(OpHope, i+1) // i+1 is the id of the test case
		}
```

```Go
	add := func(x,y) {
		x
	} hope {
		-100, 100 -> 0
		0,0 -> 0
		0, 1 -> 1
		fuzzing 5
	}
	add(1,0)
```
```Go
	fib := func(n) {
		if n <= 2{
			n
		} else {
			fib(n-1) + fib(n-2)
		}
	} hope {
		1 -> 1
		2 -> 2
		10 -> 89
	}
```
Just gives warning information and wouldn't influence the main execution
```
want 0, got -100 in the 1-th test case
want 1, got 0 in the 3-th test case
============= final result ==========
1
```
