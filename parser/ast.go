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
	return reflect.TypeOf(i).String() + " " + i.Key
}

// ==================== IntegerLiteral
type IntegerLiteral struct {
	Key int
}

func (n IntegerLiteral) String() string {
	return reflect.TypeOf(n).String() + " " + strconv.Itoa(n.Key)
}

//	=============== String
type StringLiteral struct {
	Key string
}

func (s StringLiteral) String() string {
	return reflect.TypeOf(s).String() + " " + s.Key
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

// ================= array
type ArrayLiteral struct {
	Elements []Expression
}

func (a ArrayLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(reflect.TypeOf(a).String() + " ")
	out.WriteString(LBRACE)
	for i, e := range a.Elements {
		if i != 0 {
			out.WriteString(COMMA)
		}
		out.WriteString(e.String())

	}
	out.WriteString(RBRACE)
	return out.String()
}

// =========== Index Expression
type IndexExpression struct {
	Left  Expression
	Index Expression
}

func (ie IndexExpression) String() string {
	return reflect.TypeOf(ie).String() + " " + ie.Left.String() + LBRACKET + ie.Index.String() + RBRACKET
}

// ====== function
type FunctionLiteral struct {
	Parameters []*IdentifierLiteral
	Execute    *BlockExpression
}

func (f FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(reflect.TypeOf(f).String())

	out.WriteString(LPAREN)

	paras := []string{}
	for _, p := range f.Parameters {
		paras = append(paras, p.String())
	}
	out.WriteString(strings.Join(paras, ","))

	out.WriteString(RPAREN)
	out.WriteString(f.Execute.String())

	return out.String()

}

// =========== Call Expression
type CallExpression struct {
	Function  Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(reflect.TypeOf(ce).String() + " " + ce.Function.String())
	out.WriteString(LPAREN)

	paras := []string{}
	for _, arg := range ce.Arguments {
		paras = append(paras, arg.String())
	}
	out.WriteString(strings.Join(paras, ","))

	out.WriteString(RPAREN)
	return out.String()
}

// ================== Unary Expression
type UnaryExpression struct {
	Operator string
	Right    Expression
}

func (ue UnaryExpression) String() string {
	return reflect.TypeOf(ue).String() + " " + ue.Operator + ue.Right.String()
}

// ================== Infix Expression
type InfixExpression struct {
	Operator    string
	Left, Right Expression
}

func (ie InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(reflect.TypeOf(ie).String() + " ")
	out.WriteString(LPAREN)
	out.WriteString(ie.Left.String() + ie.Operator + ie.Right.String())
	out.WriteString(RPAREN)

	return out.String()
}

// ================== AssignStatement
type AssignExpression struct {
	Ident *IdentifierLiteral
	Expr  Expression
}

func (ae AssignExpression) String() string {
	return reflect.TypeOf(ae).String() + " " + "(" + ae.Ident.String() + "=" + ae.Expr.String() + ")"
}

// define
type DefineExpression struct {
	Ident *IdentifierLiteral
	Expr  Expression
}

func (de DefineExpression) String() string {
	return reflect.TypeOf(de).String() + " " + "(" + de.Ident.String() + ":=" + de.Expr.String() + ")"
}

// ================= If Statement
type IfExpression struct {
	conditions []Expression
	executes   []*BlockExpression
}

func (ie *IfExpression) addPair(cnd Expression, block *BlockExpression) {
	ie.conditions = append(ie.conditions, cnd)
	ie.executes = append(ie.executes, block)
}

func (ie IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString(reflect.TypeOf(ie).String() + " " + "(")

	for i := range ie.conditions {
		if i == 0 {
			out.WriteString("if ")
		} else {
			out.WriteString("else if ")
		}
		out.WriteString(ie.conditions[i].String())
		out.WriteString(ie.executes[i].String())
	}
	out.WriteString(")")

	return out.String()
}

// ================= Ternary Statement
type TernaryExpression struct {
	condition   Expression
	left, right Expression
}

func (ts TernaryExpression) String() string {
	return reflect.TypeOf(ts).String() + " " + "(" + ts.condition.String() + " ? " + ts.left.String() + ":" + ts.right.String() + ")"
}

// ================= While Statement
type WhileExpression struct {
	Condition Expression
	Execute   *BlockExpression
}

func (we WhileExpression) String() string {
	return reflect.TypeOf(we).String() + " " + we.Condition.String() + we.Execute.String()
}

type BlockExpression struct {
	Statements []Statement
}

func (be BlockExpression) String() string {
	var out bytes.Buffer

	out.WriteString(LBRACE)
	for _, s := range be.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	out.WriteString(RBRACE)
	return out.String()
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
