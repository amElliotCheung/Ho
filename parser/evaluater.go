package interpreter

import (
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
		// log.Printf("receive item from channel ===========>> %T  , %v", node, node)
		result = e.evalStatement(node)
		debugs = append(debugs, result)
	}
	// show all results
	for i, item := range debugs {
		// log.Println("-- ", i+1, item.Type(), item.Inspect(), "--")
		log.Println(i, item)

	}
	return result // the last statement
}

func (e *Evaluater) evalStatement(node ASTNode) Object {
	var result Object
	switch node.(type) {
	case *IfStatement:
		log.Println("evalStatement ---> if", node)
		result = e.evalIf(node)
	case *WhileStatement:
		log.Println("evalStatement ---> While", node)
		result = e.evalWhile(node)
	default:
		log.Println("evalStatement ---> default expression", node)
		result = e.evalExpr(node)
		// case *TernaryStatement:
		// 	log.Println("evaluater ---> Ternary")
		// 	result = e.evalTernary(node)
	}
	return result
}

func (e *Evaluater) evalExpr(node ASTNode) Object {
	// if ternary
	switch node.(type) {
	case *TernaryStatement:
		log.Println("evalExpr ----> it's ternary", node)
		return e.evalTernary(node)
	case *ASTLeaf:
		log.Println("evalExpr ----> it's leaf", node)
		leaf, _ := node.(*ASTLeaf)
		if leaf.IsNumber() {
			return &Integer{Value: leaf.GetNumber()}
		}
		if leaf.IsIdentifier() {
			obj := e.environment.Get(leaf.GetText())
			if obj.Type() == "INTEGER" {
				return &Integer{Value: obj.GetValue()}
			}
		}
	case *Expression:
		log.Println("evalExpr ----> it's expression", node)
		if node.numChildren() == 2 {
			// log.Println("evaluater ----> it's binaryExpression")
			// log.Println(node, fmt.Sprintf("%T", node))
			left := e.evalExpr(node.child(0))
			right := e.evalExpr(node.child(1))
			switch node.(*Expression).operator() {
			case "=":
				name := node.child(0).(*ASTLeaf).GetText()
				e.environment.Set(name, right)
				log.Println("evalExpr ------> set a variable.", name)
				return right
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
			case ">=":
				isSatisfied := 0
				if left.GetValue() >= right.GetValue() {
					isSatisfied = 1
				}
				return &Integer{Value: isSatisfied}
			case "<=":
				isSatisfied := 0
				if left.GetValue() <= right.GetValue() {
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
		}
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

func (e *Evaluater) evalTernary(node ASTNode) Object {
	if e.isTure(node.child(0)) {
		return e.evalExpr(node.child(1))
	}
	return e.evalExpr(node.child(2))
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
