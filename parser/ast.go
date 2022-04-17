package parser

import "bytes"

type Program interface {
}
type Statement interface {
}

// ==================== ASTNode Interface
type ASTNode interface {
	child(int) ASTNode
	numChildren() int
	children() []ASTNode
	string() string
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
	return []ASTNode{}
}

func (leaf ASTLeaf) string() string {
	return leaf.Token.GetText()
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

func (list ASTList) string() string {
	var out bytes.Buffer

	out.WriteString("(")
	for _, node := range list.nodes {
		out.WriteString(node.string() + `\n`)
	}
	out.WriteString(")")

	return out.String()
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
	return ""
}
func (be BinaryExpr) right() ASTNode {
	return be.child(2)
}
