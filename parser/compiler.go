package interpreter

import (
	"encoding/binary"
	"fmt"
	"log"
)

type Instructions []byte

type CompilationScope struct {
	instructions Instructions

	// to back patch
	jumpPos        []int
	jumpIfFalsePos int
}

type Compiler struct {
	scopes []CompilationScope

	constants   []Object
	symbolTable *SymbalTable

	// for convience, when generate byte code
	operator2code map[string]Opcode
}

func NewCompiler() *Compiler {
	mainScope := CompilationScope{
		instructions:   make(Instructions, 0),
		jumpPos:        make([]int, 0),
		jumpIfFalsePos: 0,
	}
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

	return &Compiler{
		scopes:        []CompilationScope{mainScope},
		constants:     make([]Object, 0),
		symbolTable:   NewSymbolTable(),
		operator2code: operator2code,
	}
}
func (c *Compiler) Compile(node ASTNode) error {
	switch node := node.(type) {

	case *Program:
		for _, stmt := range node.Statements {
			c.Compile(stmt)
		}

	case *BlockExpression:
		for _, stmt := range node.Statements {
			c.Compile(stmt)
		}

	case *FunctionLiteral:
		c.enterScope()

		for _, para := range node.Parameters {
			c.addVariable(para.Key)
		}

		c.Compile(node.Execute)
		c.emit(OpReturnValue)
		log.Println("compiler functionliteral ---->", c.currentInstructions(), c.symbolTable.size, len(node.Parameters))
		// compiledFn := &CompiledFunction{ is wrong!!!
		compiledFn := CompiledFunction{
			Instructions: c.currentInstructions(),
			NumLocals:    c.symbolTable.size,
			NumParas:     len(node.Parameters),
		}
		log.Println("compiler functionliteral ---->", compiledFn)
		c.leaveScope()
		idx := c.addConstant(&compiledFn)
		c.emit(OpConstant, idx)
		// return nil

	case *CallExpression:
		c.Compile(node.Function)
		for _, para := range node.Arguments {
			c.Compile(para)
		}
		c.emit(OpCall, len(node.Arguments))

	case *DefineExpression:
		symbol := c.addVariable(node.Ident.Key) // return symbol

		c.Compile(node.Expr)
		if symbol.Scope == GlobalScope {
			c.emit(OpSetGlobal, symbol.Index)
		} else if symbol.Scope == LocalScope {
			c.emit(OpSetLocal, symbol.Index)
		}

	case *AssignExpression:
		c.Compile(node.Expr)
		symbol := c.getVariable(node.Ident.Key) // return symbol
		if symbol.Scope == GlobalScope {
			c.emit(OpSetGlobal, symbol.Index)
		} else if symbol.Scope == LocalScope {
			c.emit(OpSetLocal, symbol.Index)
		}

	case *IfExpression:
		for i, cnd := range node.conditions {

			c.Compile(cnd)
			c.occupy(OpJumpIfFalse)
			c.Compile(node.executes[i])
			if i == len(node.conditions)-1 { // the last block
				c.backPatch(OpJump, len(c.currentInstructions()))
			} else {
				c.occupy(OpJump)
			}
			c.backPatch(OpJumpIfFalse, len(c.currentInstructions()))

		}
	case *WhileExpression:
		cndIdx := len(c.currentInstructions())
		c.Compile(node.Condition)
		c.occupy(OpJumpIfFalse)
		c.Compile(node.Execute)
		c.emit(OpJump, cndIdx)
		c.backPatch(OpJumpIfFalse, len(c.currentInstructions()))

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
		// fmt.Println("integerLiteral debug: ", obj, idx)
		c.emit(OpConstant, idx)

	case *IdentifierLiteral:
		symbol := c.getVariable(node.Key)
		if symbol.Scope == GlobalScope {
			c.emit(OpGetGlobal, symbol.Index)
		} else if symbol.Scope == LocalScope {
			c.emit(OpGetLocal, symbol.Index)
		}

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
	case OpGetLocal, OpSetLocal:
		operand := byte(operands[0])
		ins = append(ins, operand)
	case OpCall:
		operand := byte(operands[0])
		ins = append(ins, operand)
	// no-operand opcode
	case OpPop:
	case OpAdd, OpSub, OpMult, OpDiv, OpMod:
	case OpMinus, OpBang:
	case OpReturnValue:
	}
	// add it to the list
	c.scopes[len(c.scopes)-1].instructions = append(c.scopes[len(c.scopes)-1].instructions, ins...)

	// fmt.Println("emit debug", c.scopes[len(c.scopes)-1].instructions)
	// fmt.Println("scope = ", len(c.scopes)-1)
	// fmt.Println("constants = ", c.constants)
	fmt.Println()
}

// it reserves operand for an op, particular, Jump
// this space will be filled later
func (c *Compiler) occupy(op Opcode) {
	scp := c.scopes[len(c.scopes)-1]

	switch op {
	case OpJump:
		scp.instructions = append(scp.instructions, byte(OpJump))
		scp.instructions = append(scp.instructions, make([]byte, 2)...)
		scp.jumpPos = append(scp.jumpPos, len(scp.instructions)-2)
	case OpJumpIfFalse:
		scp.instructions = append(scp.instructions, byte(OpJumpIfFalse))
		scp.instructions = append(scp.instructions, make([]byte, 2)...)
		scp.jumpIfFalsePos = len(scp.instructions) - 2
	}
	// important
	// slice must be assigned back
	c.scopes[len(c.scopes)-1] = scp
}
func (c *Compiler) backPatch(op Opcode, operands ...int) {
	scp := c.scopes[len(c.scopes)-1]
	switch op {
	case OpJump:
		operand := uint16(operands[0])
		for _, pos := range scp.jumpPos {
			binary.BigEndian.PutUint16(scp.instructions[pos:], operand)
		}
		// clear
		scp.jumpPos = scp.jumpPos[:0]
	case OpJumpIfFalse:
		operand := uint16(operands[0])
		binary.BigEndian.PutUint16(scp.instructions[scp.jumpIfFalsePos:], operand)
	}
	c.scopes[len(c.scopes)-1] = scp
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:   make(Instructions, 0),
		jumpPos:        make([]int, 0),
		jumpIfFalsePos: 0,
	}
	c.scopes = append(c.scopes, scope)

	c.symbolTable = NewEnclosedSymbalTable(c.symbolTable)
}

func (c *Compiler) leaveScope() {
	c.show()
	c.scopes = c.scopes[:len(c.scopes)-1]

	c.symbolTable = c.symbolTable.outer
}

func (c *Compiler) addVariable(name string) *Symbol {
	symbol := c.symbolTable.Define(name)
	return symbol
}

func (c *Compiler) getVariable(name string) *Symbol {
	symbol, ok := c.symbolTable.Resolve(name)
	if !ok {
		panic("viriable undefined!")
	}
	return symbol
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
		instructions: c.currentInstructions(),
		constants:    c.constants,
	}
}

func (c *Compiler) currentInstructions() Instructions {
	return c.scopes[len(c.scopes)-1].instructions
}

// debug
func (c *Compiler) show() {
	pc := 0
	log.Println("compiler.show --- > scope=", len(c.scopes)-1)
	for pc < len(c.currentInstructions()) {
		log.Println("compiler --- > pc = ", pc)
		op := Opcode(c.currentInstructions()[pc])
		switch op {
		case OpConstant:
			log.Println("compiler --- > constant ", c.currentInstructions()[pc:pc+3])
			pc += 3
		case OpPop:
			log.Println("compiler --- > pop")
			pc++

		case OpAdd, OpSub, OpMult, OpDiv, OpMod, OpLt, OpGt, OpLte, OpGte, OpEq, OpNeq:
			log.Println("compiler --- > add sub lt ...")
			pc++

		case OpMinus:
			log.Println("compiler --- > minus")
			pc++

		case OpBang:
			log.Println("compiler --- > bang")
			pc++

		case OpJump:
			log.Println("compiler --- > jump  ", c.currentInstructions()[pc:pc+3])
			pc = int(binary.BigEndian.Uint16(c.currentInstructions()[pc+1:]))

		case OpJumpIfFalse:

			log.Println("compiler --- > jumpIfFalse", c.currentInstructions()[pc:pc+3])
			pc += 3

		case OpSetGlobal:
			log.Println("compiler --- > setGlobal  ", c.currentInstructions()[pc:pc+3])
			pc += 3

		case OpGetGlobal:
			log.Println("compiler --- > getGlobal  ", c.currentInstructions()[pc:pc+3])
			pc += 3

		case OpSetLocal:
			log.Println("compiler --- > setLocal  ", c.currentInstructions()[pc:pc+2])
			pc += 2

		case OpGetLocal:
			log.Println("compiler --- > getLocal  ", c.currentInstructions()[pc:pc+2])
			pc += 2

		case OpCall:
			log.Println("compiler --- > call  ", c.currentInstructions()[pc:pc+2])
			pc += 1
		default:
			return
		}
		// debug
	}
}
