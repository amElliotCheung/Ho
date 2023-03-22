package interpreter

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
)

//
//	Interfaces
//

// ==================== ASTNode Interface
type ASTNode interface {
	String() string
	Type() string
}

// ==================== Statement Interface
type Statement interface {
	ASTNode
}

// ====================  Expression Interface
type Expression interface {
	ASTNode
}

//
// data types
//
type IdentifierLiteral struct {
	Key string
}

func (i IdentifierLiteral) String() string {
	return i.Key
}
func (i IdentifierLiteral) Type() string {
	return "IdentifierLiteral"
}

// ==================== IntegerLiteral
type IntegerLiteral struct {
	Key int
}

func (n IntegerLiteral) String() string {
	return strconv.Itoa(n.Key)
}
func (n IntegerLiteral) Type() string {
	return "IntegerLiteral"
}

//	=============== String
type StringLiteral struct {
	Key string
}

func (s StringLiteral) String() string {
	return s.Key
}
func (s StringLiteral) Type() string {
	return "StringLiteral"
}

// ================= boolean
type BooleanLiteral struct {
	Key bool
}

func (b BooleanLiteral) String() string {
	if b.Key {
		return "true"
	}
	return "false"
}
func (b BooleanLiteral) Type() string {

	return "BooleanLiteral"
}

// ================= array
type ArrayLiteral struct {
	Elements []Expression
}

func (a ArrayLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(LBRACKET)
	for i, e := range a.Elements {
		if i != 0 {
			out.WriteString(COMMA)
		}
		out.WriteString(e.String())

	}
	out.WriteString(RBRACKET)
	return out.String()
}
func (a ArrayLiteral) Type() string {
	return "ArrayLiteral"
}

// =========== Index Expression
type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (ie IndexExpression) String() string {
	return ie.Left.String() + LBRACKET + ie.Index.String() + RBRACKET
}
func (ie IndexExpression) Type() string {
	return "IndexExpression"
}

// ====== function
type FunctionLiteral struct {
	Parameters []*IdentifierLiteral
	ParaTypes  []string
	Execute    *BlockExpression
	Hopes      *HopeBlock
}

func (f FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("func")
	out.WriteString(LPAREN)

	paras := []string{}
	for _, p := range f.Parameters {
		paras = append(paras, p.String())
	}
	out.WriteString(strings.Join(paras, ","))

	out.WriteString(RPAREN)
	out.WriteString(f.Execute.String())

	if f.Hopes != nil {
		out.WriteString(" ")
		out.WriteString(f.Hopes.String())
	}

	return out.String()

}
func (f FunctionLiteral) Type() string {
	return "FunctionLiteral"
}

type HopeBlock struct {
	HopeExpressions []HopeExpression
	NFuzzing        *IntegerLiteral
}

func (hb HopeBlock) String() string {
	var out bytes.Buffer
	out.WriteString("hope { \n")
	for _, expr := range hb.HopeExpressions {
		out.WriteString(expr.String())
		out.WriteString("\n")
	}
	if hb.NFuzzing != nil {
		out.WriteString(FUZZING + " " + hb.NFuzzing.String())
	}
	out.WriteString("\n" + RBRACE)

	return out.String()
}
func (hb HopeBlock) Type() string {
	return "HopeBlock"
}

type HopeExpression struct {
	Parameters []Expression
	Expected   Expression
}

func (hp HopeExpression) String() string {
	var out bytes.Buffer

	paras := []string{}
	for _, p := range hp.Parameters {
		paras = append(paras, p.String())
	}
	out.WriteString(strings.Join(paras, ","))

	out.WriteString(" -> " + hp.Expected.String())

	return out.String()
}
func (hp HopeExpression) Type() string {
	return "HopeExpression"
}

// =========== Call Expression
type CallExpression struct {
	Function  Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString(LPAREN)

	paras := []string{}
	for _, arg := range ce.Arguments {
		paras = append(paras, arg.String())
	}
	out.WriteString(strings.Join(paras, ","))

	out.WriteString(RPAREN)

	return out.String()
}
func (ce CallExpression) Type() string {
	return "CallExpression"
}

// ================== Unary Expression
type UnaryExpression struct {
	Operator string
	Right    Expression
}

func (ue UnaryExpression) String() string {
	return ue.Operator + ue.Right.String()
}
func (ue UnaryExpression) Type() string {
	return "UnaryExpression"
}

// ================== Infix Expression
type InfixExpression struct {
	Operator    string
	Left, Right Expression
}

func (ie InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.Left.String() + ie.Operator + ie.Right.String())

	return out.String()
}
func (ie InfixExpression) Type() string {
	return "InfixExpression"
}

// ================== AssignStatement
type AssignExpression struct {
	Ident *IdentifierLiteral
	Expr  Expression
}

func (ae AssignExpression) String() string {
	return ae.Ident.String() + " = " + ae.Expr.String()
}
func (ae AssignExpression) Type() string {
	return "AssignExpression"
}

// define
type DefineExpression struct {
	Ident *IdentifierLiteral
	Expr  Expression
}

func (de DefineExpression) String() string {
	return de.Ident.String() + ":=" + de.Expr.String()
}
func (de DefineExpression) Type() string {
	return "DefineExpression"
}

// ================= If Statement
type IfExpression struct {
	conditions []Expression
	executes   []*BlockExpression
}

// must be a pointer!
func (ie *IfExpression) addPair(cnd Expression, block *BlockExpression) {
	ie.conditions = append(ie.conditions, cnd)
	ie.executes = append(ie.executes, block)
}

func (ie IfExpression) String() string {
	var out bytes.Buffer

	for i := range ie.conditions {
		if i == 0 {
			out.WriteString("if ")
		} else {
			out.WriteString("else if ")
		}
		out.WriteString(ie.conditions[i].String())
		out.WriteString(ie.executes[i].String())
	}

	return out.String()
}
func (ie IfExpression) Type() string {
	return "IfExpression"
}

// ================= Ternary Statement
type TernaryExpression struct {
	condition   Expression
	left, right Expression
}

func (ts TernaryExpression) String() string {
	return ts.condition.String() + " ? " + ts.left.String() + " : " + ts.right.String() + ")"
}
func (ts TernaryExpression) Type() string {
	return "TernaryExpression"
}

// ================= While Statement
type WhileExpression struct {
	Condition Expression
	Execute   *BlockExpression
}

func (we WhileExpression) String() string {
	return "while " + we.Condition.String() + we.Execute.String()
}
func (we WhileExpression) Type() string {
	return "WhileExpression"
}

type BlockExpression struct {
	Statements []Statement
}

func (be BlockExpression) String() string {
	var out bytes.Buffer

	out.WriteString(LBRACE + "\n")
	for _, s := range be.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	out.WriteString(RBRACE)
	return out.String()
}
func (be BlockExpression) Type() string {
	return "BlockExpression"
}

// ==================== Program
type Program struct {
	Statements []Statement
}

func (p Program) String() string {
	var out bytes.Buffer

	out.WriteString(reflect.TypeOf(p).String() + "\n")
	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (p Program) Type() string {
	return "Program"
}
