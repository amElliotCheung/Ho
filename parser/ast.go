package interpreter

import (
	"bytes"
	"fmt"
)

// ==================== ASTNode Interface
type ASTNode interface {
	child(int) ASTNode
	numChildren() int
	children() []ASTNode
	String() string
	addChild(ASTNode)
}

// ==================== ASTLeaf
type ASTLeaf struct {
	Token
}

func (leaf ASTLeaf) child(i int) ASTNode {
	return nil
}
func (leaf ASTLeaf) numChildren() int {
	return 0
}
func (leaf ASTLeaf) children() []ASTNode {
	return nil
}

func (leaf ASTLeaf) String() string {
	return leaf.Token.GetText()
}

func (leaf ASTLeaf) addChild(node ASTNode) {
	panic("a leaf can't have child")
}

// ==================== NumberLiteral
type NumberLiteral struct {
	ASTLeaf
}

func (number NumberLiteral) value() int {
	return number.ASTLeaf.GetNumber()
}

// ==================== ASTList
type ASTList struct {
	Token
	nodes []ASTNode
}

func (list ASTList) child(i int) ASTNode {
	if i < len(list.nodes) {
		return list.nodes[i]
	}
	return nil
}
func (list ASTList) numChildren() int {
	return len(list.nodes)
}
func (list ASTList) children() []ASTNode {
	return list.nodes
}

func (list ASTList) String() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("I have %d children ==> ", len(list.nodes)))
	out.WriteString("(")

	for _, node := range list.nodes {
		out.WriteString(node.String() + " ")
	}
	out.WriteString(")")

	return out.String()
}

func (l *ASTList) addChild(node ASTNode) {
	l.nodes = append(l.nodes, node)
}

// ================== Expression
type Expression struct {
	ASTList
}

func (be Expression) left() ASTNode {
	return be.child(0)
}
func (be Expression) operator() string {
	// later
	return be.Token.GetText()
}
func (be Expression) right() ASTNode {
	return be.child(1)
}
func (be Expression) String() string {
	return "(" + be.child(0).String() + be.Token.GetText() + be.child(1).String() + ")"
}

// ================== AssignStatement
type AssignStatement struct {
	ASTList
}

func (as AssignStatement) ident() ASTNode {
	return as.child(0)
}

func (as AssignStatement) value() ASTNode {
	return as.child(1)
}

func (as AssignStatement) String() string {
	return "(" + as.child(0).String() + "=" + as.child(1).String() + ")"
}

// ================= If Statement
type IfStatement struct {
	ASTList
}

func (is IfStatement) condition(i int) ASTNode {
	return is.child(2 * i)
}

func (is IfStatement) block(i int) ASTNode {
	return is.child(2*i + 1)
}

func (is IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("(")

	for i := 0; i < len(is.nodes); i += 2 {
		if i == 0 {
			out.WriteString("if ")
		} else {
			out.WriteString("else if ")
		}
		out.WriteString(is.child(i).String())
		out.WriteString("then " + is.child(i+1).String())
	}
	out.WriteString(")")

	return out.String()
}

// ================= Ternary Statement
type TernaryStatement struct {
	ASTList
}

// ================= While Statement
type WhileStatement struct {
	ASTList
}

func (ws WhileStatement) condition() ASTNode {
	return ws.child(0)
}

func (ws WhileStatement) block() ASTNode {
	return ws.child(1)
}

func (ws WhileStatement) String() string {
	return "(" + ws.Token.GetText() + " " + ws.child(0).String() + " do " + ws.child(1).String() + ")"
}

// ================= func Statement
type FuncStatement struct {
	parameters []*IdToken
	ASTList
}

func (ws FuncStatement) condition() ASTNode {
	return ws.child(0)
}

func (ws FuncStatement) block() ASTNode {
	return ws.child(1)
}

func (ws FuncStatement) String() string {
	return "(" + ws.Token.GetText() + " " + ws.child(0).String() + " do " + ws.child(1).String() + ")"
}
