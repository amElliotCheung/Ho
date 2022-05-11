package interpreter

import (
	"encoding/binary"
	"fmt"
	"log"
)

type Instructions []byte
type Opcode byte

const (
	OpConstant Opcode = iota
	OpPop

	// operators
	OpAdd
	OpSub
	OpMult
	OpDiv
	OpMod
	OpMinus
	OpBang
	OpLt
	OpGt
	OpLte
	OpGte

	OpEq
	OpNeq

	// Jump
	OpJump
	OpJumpIfFalse

	// variables
	OpGetGlobal
	OpSetGlobal
)

const (
	GlobalScope = "Global"
)

type Symbol struct {
	name, scope string
	index       int
}

type SymbalTable struct {
	store map[string]*Symbol
	size  int
}

type Compiler struct {
	instructions Instructions
	constants    []Object
	symbolTable  SymbalTable
	// for convience, when generate byte code
	operator2code map[string]Opcode
	// to back patch
	jumpPos        []int
	jumpIfFalsePos int
}

func NewCompiler() *Compiler {
	operator2code := map[string]Opcode{PLUS: OpAdd,
		MINUS:    OpSub,
		ASTERISK: OpMult,
		SLASH:    OpDiv,
		MOD:      OpMod,
		BANG:     OpBang,
		LT:       OpLt,
		GT:       OpGt,
		LTE:      OpLte,
		GTE:      OpGte,

		EQ:  OpEq,
		NEQ: OpNeq,
	}

	symbolTable := SymbalTable{
		store: make(map[string]*Symbol),
	}
	return &Compiler{
		instructions:   make(Instructions, 0),
		constants:      make([]Object, 0),
		symbolTable:    symbolTable,
		operator2code:  operator2code,
		jumpPos:        make([]int, 0),
		jumpIfFalsePos: 0,
	}
}
func (c *Compiler) Compile(node ASTNode) error {
	switch node := node.(type) {
	case *Program:
		for _, stmt := range node.Statements {
			c.Compile(stmt)
		}
	case *IfExpression:
		for i, cnd := range node.conditions {

			c.Compile(cnd)
			c.occupy(OpJumpIfFalse)
			c.Compile(node.executes[i])
			if i == len(node.conditions)-1 { // the last block
				c.backPatch(OpJump, len(c.instructions))
			} else {
				c.occupy(OpJump)
			}
			c.backPatch(OpJumpIfFalse, len(c.instructions))

		}

	case *AssignExpression:
		c.Compile(node.Expr)
		idx := c.addVariable(node.Ident.Key) // return index
		c.emit(OpSetGlobal, idx)

	case *WhileExpression:
		cndIdx := len(c.instructions)
		c.Compile(node.Condition)
		c.occupy(OpJumpIfFalse)
		c.Compile(node.Execute)
		c.emit(OpJump, cndIdx)
		c.backPatch(OpJumpIfFalse, len(c.instructions))

	case *BlockExpression:
		for _, stmt := range node.Statements {
			c.Compile(stmt)
		}

	case *InfixExpression:
		c.Compile(node.Left)
		c.Compile(node.Right)
		code, ok := c.operator2code[node.Operator]
		if !ok {
			return fmt.Errorf("illegal operator infix expression")
		}
		c.emit(code)

	case *UnaryExpression:
		c.Compile(node.Right)

		switch node.Operator {
		case BANG:
			c.emit(OpBang)
		case MINUS:
			c.emit(OpMinus)
		default:
			return fmt.Errorf("illegal operator unary expression")
		}

	case *IntegerLiteral:
		obj := &Integer{Value: node.Key}
		idx := c.addConstant(obj)
		c.emit(OpConstant, idx)

	case *IdentifierLiteral:
		idx := c.getVariableIndex(node.Key)
		c.emit(OpGetGlobal, idx)

	case *BooleanLiteral:
		obj := &Boolean{Value: node.Key}
		idx := c.addConstant(obj)
		c.emit(OpConstant, idx)

	case *StringLiteral:
		obj := &String{Value: node.Key}
		idx := c.addConstant(obj)
		c.emit(OpConstant, idx)
	}
	return nil
}

func (c *Compiler) emit(op Opcode, operands ...int) {
	// make an instruction
	ins := Instructions{0: byte(op)}
	switch op {
	case OpConstant, OpSetGlobal, OpGetGlobal, OpJump: // only one width-2 operand, the constant index
		operand := uint16(operands[0])
		ins = append(ins, byte(operand>>8))
		ins = append(ins, byte(operand))

	// no-operand opcode
	// simply write them down
	case OpPop:
	case OpAdd, OpSub, OpMult, OpDiv, OpMod:
	case OpMinus, OpBang:
	}
	// add it to the list
	c.instructions = append(c.instructions, ins...)
}

// it reserves operand for an op, particular, Jump
// this space will be filled later
func (c *Compiler) occupy(op Opcode) {
	switch op {
	case OpJump:
		c.instructions = append(c.instructions, byte(OpJump))
		c.instructions = append(c.instructions, make([]byte, 2)...)
		c.jumpPos = append(c.jumpPos, len(c.instructions)-2)
	case OpJumpIfFalse:
		c.instructions = append(c.instructions, byte(OpJumpIfFalse))
		c.instructions = append(c.instructions, make([]byte, 2)...)
		c.jumpIfFalsePos = len(c.instructions) - 2
	}
}
func (c *Compiler) backPatch(op Opcode, operands ...int) {
	switch op {
	case OpJump:
		operand := uint16(operands[0])
		for _, pos := range c.jumpPos {
			binary.BigEndian.PutUint16(c.instructions[pos:], operand)
		}
		// clear
		c.jumpPos = c.jumpPos[:0]
	case OpJumpIfFalse:
		operand := uint16(operands[0])
		binary.BigEndian.PutUint16(c.instructions[c.jumpIfFalsePos:], operand)
	}
}
func (c *Compiler) addVariable(name string) int {
	if sb, ok := c.symbolTable.store[name]; ok { // exist
		return sb.index
	}
	idx := c.symbolTable.size
	c.symbolTable.store[name] = &Symbol{
		name:  name,
		scope: GlobalScope,
		index: idx,
	}
	c.symbolTable.size++
	return idx
}
func (c *Compiler) getVariableIndex(name string) int {
	if sb, ok := c.symbolTable.store[name]; ok { // exist
		return sb.index
	}
	panic("viriable undefined!")
}

func (c *Compiler) addConstant(obj Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// output result to virual machine
type Bytecode struct {
	instructions Instructions
	constants    []Object
}

func (c *Compiler) bytecode() Bytecode {
	return Bytecode{
		instructions: c.instructions,
		constants:    c.constants,
	}
}

// debug
func (c *Compiler) show() {
	pc := 0
	for pc < len(c.instructions) {
		log.Println("pc = ", pc)
		op := Opcode(c.instructions[pc])
		switch op {
		case OpConstant:
			idx := binary.BigEndian.Uint16(c.instructions[pc+1:])
			log.Println("constant ", c.constants[idx])
			pc += 3
		case OpPop:
			log.Println("pop")
			pc++

		case OpAdd, OpSub, OpMult, OpDiv, OpMod, OpLt, OpGt, OpLte, OpGte, OpEq, OpNeq:
			log.Println("add sub lt ...")
			pc++

		case OpMinus:
			log.Println("minus")
			pc++

		case OpBang:
			log.Println("bang")
			pc++

		case OpJump:
			idx := binary.BigEndian.Uint16(c.instructions[pc+1:])
			log.Println("jump  ", idx)
			pc = int(binary.BigEndian.Uint16(c.instructions[pc+1:]))

		case OpJumpIfFalse:

			log.Println("jumpIfFalse", int(binary.BigEndian.Uint16(c.instructions[pc+1:])))
			pc += 3

		case OpSetGlobal:
			idx := int(binary.BigEndian.Uint16(c.instructions[pc+1:]))
			log.Println("setGlobal  ", idx)
			pc += 3

		case OpGetGlobal:
			idx := int(binary.BigEndian.Uint16(c.instructions[pc+1:]))
			log.Println("getGlobal  ", idx)
			pc += 3
		}
		// debug
	}
}
