package interpreter

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

	OpGetLocal
	OpSetLocal

	// function
	OpCall
	OpReturnValue
)

const (
	GlobalScope = "Global"
	LocalScope  = "Local"
)
