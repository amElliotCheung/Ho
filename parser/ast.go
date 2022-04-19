package parser

import (
	"bytes"
	"fmt"
)

type Program interface {
}
type Statement interface {
}

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

// ==================== Name
type Name struct {
	ASTLeaf
}

func (name Name) name() string {
	return name.ASTLeaf.GetText()
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

	out.WriteString(fmt.Sprintf("i have %d children", len(list.nodes)))
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

// ================== BinaryExpr
type BinaryExpr struct {
	ASTList
}

func (be BinaryExpr) left() ASTNode {
	return be.child(0)
}
func (be BinaryExpr) operator() string {
	// later
	return be.Token.GetType()
}
func (be BinaryExpr) right() ASTNode {
	return be.child(1)
}
func (be BinaryExpr) String() string {
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
	return "(" + as.child(0).String() + as.Token.GetText() + as.child(1).String() + ")"
}
