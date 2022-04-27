package interpreter

import (
	"fmt"
	"log"
)

type Evaluater struct {
	c           chan ASTNode
	environment *Environment
}

func NewEvaluater(c chan ASTNode) *Evaluater {
	return &Evaluater{c: c,
		environment: NewEnvironment()}
}
func (e *Evaluater) Eval() Object {
	var result Object
	debugs := make([]Object, 0)
	for node := range e.c {
		log.Println("receive item from channel ===========>>>> ", node)
		result = e.evalStatement(node)
		debugs = append(debugs, result)
	}
	// show all results
	for i, item := range debugs {
		log.Println("--  ", i+1, item.Type(), item.Inspect(), "--")

	}
	return result // the last statement
}

func (e *Evaluater) evalStatement(node ASTNode) Object {
	var result Object
	switch node.(type) {
	case *ASTLeaf: // just a token. Identifier or number
		log.Println("evaluater ---> integer")
		leaf, _ := node.(*ASTLeaf)
		result = &Integer{Value: leaf.GetNumber()}
	case *Expression:
		log.Println("evaluater ---> expression")
		result = e.evalExpr(node)
	case *IfStatement:
		log.Println("evaluater ---> if")
		result = e.evalIf(node)
	case *AssignStatement:
		log.Println("evaluater ---> assign")
		result = e.evalAssign(node)
	case *WhileStatement:
		log.Println("evaluater ---> while")
		result = e.evalWhile(node)
	}
	return result
}

func (e *Evaluater) evalExpr(node ASTNode) Object {
	// if leaf

	if leaf, isLeaf := node.(*ASTLeaf); isLeaf {
		log.Println("evaluater ----> it's leaf")
		if leaf.IsNumber() {
			return &Integer{Value: leaf.GetNumber()}
		}
		if leaf.IsIdentifier() {
			obj := e.environment.Get(leaf.GetText())
			if obj.Type() == "INTEGER" {
				return &Integer{Value: obj.GetValue()}
			}
		}
	}
	// if binary expression
	log.Println("evaluater ----> it's binaryExpression")
	log.Println(node, fmt.Sprintf("%T", node))
	expr, _ := node.(*Expression)
	left := e.evalExpr(expr.left())
	right := e.evalExpr(expr.right())
	switch expr.operator() {
	case "+":
		return &Integer{Value: left.GetValue() + right.GetValue()}
	case "-":
		return &Integer{Value: left.GetValue() - right.GetValue()}
	case "*":
		return &Integer{Value: left.GetValue() * right.GetValue()}
	case "/":
		return &Integer{Value: left.GetValue() / right.GetValue()}
	case "%":
		return &Integer{Value: left.GetValue() % right.GetValue()}
	case ">":
		isSatisfied := 0
		if left.GetValue() > right.GetValue() {
			isSatisfied = 1
		}
		return &Integer{Value: isSatisfied}
	case "<":
		isSatisfied := 0
		if left.GetValue() < right.GetValue() {
			isSatisfied = 1
		}
		return &Integer{Value: isSatisfied}
	case "==":
		isSatisfied := 0
		if left.GetValue() == right.GetValue() {
			isSatisfied = 1
		}
		return &Integer{Value: isSatisfied}
	case "!=":
		isSatisfied := 0
		if left.GetValue() != right.GetValue() {
			isSatisfied = 1
		}
		return &Integer{Value: isSatisfied}
	}
	return nil
}

func (e *Evaluater) evalIf(node ASTNode) Object {
	for i := 0; i < node.numChildren(); i += 2 {
		if e.isTure(node.child(i)) {
			return e.evalBlock(node.child(i + 1))
		}
	}
	return &Integer{Value: 0}
}

func (e *Evaluater) evalBlock(node ASTNode) Object {
	var result Object
	for _, child := range node.children() {
		result = e.evalStatement(child)
	}
	return result
}

func (e *Evaluater) evalAssign(node ASTNode) Object {
	ident := node.child(0).(*ASTLeaf).GetText()
	value := e.evalExpr(node.child(1))
	e.environment.Set(ident, value)
	return value
}

func (e *Evaluater) evalWhile(node ASTNode) Object {
	var result Object
	condition := node.child(0)
	body := node.child(1)
	for e.isTure(condition) {
		result = e.evalBlock(body)
	}
	return result
}

// ==================== helper functions
func (e *Evaluater) isTure(node ASTNode) bool {
	res := e.evalExpr(node)
	return res.Type() == "INTEGER" && res.GetValue() != 0
}
