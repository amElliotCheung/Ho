package interpreter

import (
	"encoding/binary"
	"fmt"
	"log"
)

const StackSize = 2048
const VariableSize = 2048

type VM struct {
	instructions Instructions
	constants    []Object
	globals      []Object

	stack []Object
	top   int
}

func NewVM(bc Bytecode) *VM {
	return &VM{
		instructions: bc.instructions,
		constants:    bc.constants,
		stack:        make([]Object, StackSize),
		globals:      make([]Object, VariableSize),
		top:          0,
	}
}

func (vm *VM) Run() error {
	pc := 0
	iter := 0
	for pc < len(vm.instructions) {
		op := Opcode(vm.instructions[pc])
		switch op {
		case OpConstant:
			log.Println("constant")
			idx := binary.BigEndian.Uint16(vm.instructions[pc+1:])
			vm.push(vm.constants[idx])
			pc += 3
		case OpPop:
			log.Println("pop")

		case OpAdd, OpSub, OpMult, OpDiv, OpMod, OpLt, OpGt, OpLte, OpGte, OpEq, OpNeq:
			log.Println("add sub lt ...")
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
			pc++

		case OpMinus:
			log.Println("minus")
			right := vm.pop().(*Integer)
			right.Value = -right.Value
			vm.push(right)
			pc++

		case OpBang:
			log.Println("bang")
			right := vm.pop().(*Boolean)
			right.Value = !right.Value
			vm.push(right)
			pc++

		case OpJump:
			log.Println("jump")
			pc = int(binary.BigEndian.Uint16(vm.instructions[pc+1:]))

		case OpJumpIfFalse:
			log.Println("jumpIfFalse")

			if cnd := vm.pop().(*Boolean); !cnd.Value {
				log.Println("false")
				pc = int(binary.BigEndian.Uint16(vm.instructions[pc+1:]))
			} else {
				log.Println("true")
				pc += 3
			}

		case OpSetGlobal:
			log.Println("setGlobal")
			idx := int(binary.BigEndian.Uint16(vm.instructions[pc+1:]))
			vm.globals[idx] = vm.pop()
			pc += 3

		case OpGetGlobal:
			log.Println("getGlobal")
			idx := int(binary.BigEndian.Uint16(vm.instructions[pc+1:]))
			obj := vm.globals[idx]
			vm.push(obj)
			pc += 3
		}
		// debug
		iter++
		if iter > 2048*1024 {
			return fmt.Errorf("for ever")
		}
		// log.Printf("============= %dth iter pc=%d ==========", iter, pc)
		for i := 0; i < 5; i++ {
			fmt.Printf("%v , ", vm.globals[i])
		}
		// fmt.Println()
		for i := 0; i < 10; i++ {
			log.Printf("%v ", vm.stack[i])
		}
		// log.Println()
	}
	log.Printf("============= final result ==========")
	log.Println(vm.pop())
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
	if vm.top > StackSize {
		panic("stack overflow")
	}
	vm.stack[vm.top] = obj
	vm.top++
}

func (vm *VM) pop() Object {
	if vm.top == 0 {
		panic("empty stack!")
	}
	vm.top -= 1
	return vm.stack[vm.top]
}
