package interpreter

import (
	"fmt"
	"log"
)

const StackSize = 1 << 18
const VariableSize = 1 << 10

type VM struct {
	constants []Object
	globals   []Object

	frames []*Frame

	stack    []Object
	stackIdx int
}

func NewVM(bc Bytecode) *VM {
	mainFn := &CompiledFunction{Instructions: bc.instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := []*Frame{mainFrame}

	return &VM{
		constants: bc.constants,

		stackIdx: 0,
		stack:    make([]Object, StackSize),
		globals:  make([]Object, VariableSize),

		frames: frames,
	}
}

func (vm *VM) Run() error {
	// iter := 0
	for frame := vm.currentFrame(); frame.ip < len(frame.fn.Instructions); frame = vm.currentFrame() {
		op := Opcode(frame.fn.Instructions[frame.ip])
		log.Println("\n\n\n===============================")
		switch op {

		case OpConstant:
			log.Println("constant")
			idx := frame.readUint16(1)
			vm.push(vm.constants[idx])
			frame.ip += 3

		case OpPop:
			log.Println("pop")
			vm.pop()
			frame.ip++

		case OpAdd, OpSub, OpMult, OpDiv, OpMod, OpLt, OpGt, OpLte, OpGte, OpEq, OpNeq:
			log.Println("add sub lt ...", OpAdd)
			right := vm.pop()
			left := vm.pop()
			if left.Type() != right.Type() {
				log.Panicln("different type in infix expression")
			}
			switch left.Type() {
			case INTEGER_OBJ:
				vm.integerInfix(op, left, right)
			case STRING_OBJ:
				vm.stringInfix(op, left, right)
			}
			frame.ip++

		case OpMinus:
			log.Println("minus")
			right := vm.pop().(*Integer)
			right.Value = -right.Value
			vm.push(right)
			frame.ip++

		case OpBang:
			log.Println("bang")
			right := vm.pop().(*Boolean)
			right.Value = !right.Value
			vm.push(right)
			frame.ip++

		case OpJump:
			log.Println("jump")
			frame.ip = frame.readUint16(1)

		case OpJumpIfFalse:
			log.Println("jumpIfFalse")

			if cnd := vm.pop().(*Boolean); !cnd.Value {
				log.Println("false")
				frame.ip = frame.readUint16(1)

			} else {
				log.Println("true")
				frame.ip += 3
			}

		case OpSetGlobal:
			log.Println("setGlobal")
			idx := frame.readUint16(1)
			vm.globals[idx] = vm.pop()
			frame.ip += 3

		case OpGetGlobal:
			log.Println("getGlobal")
			idx := frame.readUint16(1)
			obj := vm.globals[idx]
			vm.push(obj)
			frame.ip += 3

		case OpCall:
			log.Println("call")
			numParas := frame.readUint8(1)
			fn := vm.stack[vm.stackIdx-1-numParas].(*CompiledFunction)
			nextFrame := NewFrame(fn, vm.stackIdx-numParas)
			// next Frame : important!!
			vm.stackIdx = nextFrame.bp + nextFrame.fn.NumLocals
			vm.pushFrame(nextFrame) // base pointer is current stack index
			frame.ip += 2

		case OpGetLocal:
			log.Println("get local")
			idx := frame.readUint8(1)
			obj := vm.stack[frame.bp+idx]
			vm.push(obj)
			frame.ip += 2

		case OpSetLocal:
			log.Println("set local")
			idx := frame.readUint8(1)
			obj := vm.pop()
			vm.stack[frame.bp+idx] = obj
			frame.ip += 2

		case OpReturnValue:
			log.Println("return value")
			result := vm.pop()
			f := vm.popFrame()
			vm.stack[f.bp-1] = result
			vm.stackIdx = f.bp
		}

		// debug
		// iter++
		// if iter > 2048*1024 {
		// 	return fmt.Errorf("for ever")
		// }
		// log.Printf("============= %dth iter pc=%d ==========", iter, pc)
		log.Println("vm ------- globals : ")
		for i := 0; i < 5; i++ {
			log.Printf("%v , ", vm.globals[i])
		}
		// fmt.Println()
		log.Println()
		log.Println("vm ------- stack :  (stackIdx=)", vm.stackIdx)
		for i := 0; i < vm.stackIdx; i++ {
			log.Printf("%v ", vm.stack[i])
		}
		// log.Println()
	}
	fmt.Println("============= final result ==========")
	fmt.Println(vm.stackIdx)
	fmt.Println(vm.pop())
	return nil
}

func (vm *VM) integerInfix(code Opcode, left, right Object) {
	l := left.(*Integer).Value
	r := right.(*Integer).Value

	var obj Object
	switch code {
	case OpAdd:
		obj = &Integer{Value: l + r}
	case OpSub:
		obj = &Integer{Value: l - r}
	case OpMult:
		obj = &Integer{Value: l * r}
	case OpDiv:
		obj = &Integer{Value: l / r}
	case OpMod:
		obj = &Integer{Value: l % r}
	case OpLt:
		obj = &Boolean{Value: l < r}
	case OpLte:
		obj = &Boolean{Value: l <= r}
	case OpGt:
		obj = &Boolean{Value: l > r}
	case OpGte:
		obj = &Boolean{Value: l >= r}
	case OpEq:
		obj = &Boolean{Value: l == r}
	case OpNeq:
		obj = &Boolean{Value: l != r}

	}
	vm.push(obj)
}

func (vm *VM) stringInfix(code Opcode, left, right Object) {
	l := left.(*String).Value
	r := right.(*String).Value

	var obj Object
	switch code {
	case OpAdd:
		obj = &String{Value: l + r}
	case OpEq:
		obj = &Boolean{Value: l == r}
	case OpNeq:
		obj = &Boolean{Value: l != r}

	default:
		panic("illegal operator for string")
	}
	vm.push(obj)
}

func (vm *VM) push(obj Object) {
	if vm.stackIdx > StackSize {
		panic("stack overflow")
	}
	vm.stack[vm.stackIdx] = obj
	vm.stackIdx++
}

func (vm *VM) pop() Object {
	if vm.stackIdx == 0 {
		panic("empty stack!")
	}
	vm.stackIdx--
	return vm.stack[vm.stackIdx]
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[len(vm.frames)-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames = append(vm.frames, f)
}

func (vm *VM) popFrame() *Frame {
	f := vm.currentFrame()
	vm.frames = vm.frames[:len(vm.frames)-1]
	return f

}
