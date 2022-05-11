package interpreter

import (
	"log"
)

type Evaluater struct {
	c   chan Statement
	env *Environment
}

func NewEvaluater(c chan Statement) *Evaluater {
	return &Evaluater{c: c,
		env: NewEnvironment(nil)}
}
func (e *Evaluater) Eval() Object {
	var result Object
	debugs := make([]Object, 0)
	for node := range e.c {
		log.Printf("receive item from channel ===========>> (%T, %v)", node, node)
		result = e.eval(node)
		debugs = append(debugs, result)
	}
	// show all results
	for i, item := range debugs {
		log.Println("-- ", i+1, item.Type(), item.String(), "--")
		// log.Println(i, item)

	}
	return result // the last statement
}

func (e *Evaluater) eval(node Statement) Object {
	var result Object

	switch node := node.(type) {
	case *IdentifierLiteral:
		if obj, err := e.env.Get(node.Key); err == nil {
			result = obj
		} else {
			if bt, ok := builtins[node.Key]; ok {
				result = bt
			} else {
				log.Panicln("eval IdentifierLiteral : ", err)
			}
		}
	case *IntegerLiteral:
		result = &Integer{Value: node.Key}
	case *StringLiteral:
		result = &String{Value: node.Key}
	case *BooleanLiteral:
		result = &Boolean{Value: node.Key}
	case *DefineExpression:
		result = e.evalDefine(node)
	case *AssignExpression:
		result = e.evalAssign(node)
	case *FunctionLiteral:
		// log.Println("Function literal", node)
		paras := node.Parameters
		body := node.Execute
		result = &Function{Parameters: paras, Body: body, Env: e.env}
	case *ArrayLiteral:
		result = e.evalArray(node)

	case *IndexExpression:
		log.Println("index expression")
		left := e.eval(node.Left)
		idx := e.eval(node.Index)
		result = e.evalIndex(left, idx)
	case *TernaryExpression:
		result = e.evalTernary(node)
	case *IfExpression:
		log.Println("evalStatement ---> if", node)

		result = e.evalIf(node)
	case *WhileExpression:
		log.Println("evalStatement ---> While", node)
		result = e.evalWhile(node)

	case *CallExpression:
		fn := e.eval(node.Function)         // return function object
		args := e.evalExprs(node.Arguments) // expressions
		log.Println("------- fn args ---------")
		log.Println(fn, "  ,  ", args)
		result = e.applyFunction(fn, args)

	case *UnaryExpression: // ! -
		obj := e.eval(node.Right)
		switch {
		case node.Operator == MINUS && obj.Type() == INTEGER_OBJ:
			return &Integer{Value: -obj.(*Integer).Value}

		case node.Operator == BANG && obj.Type() == BOOLEAN_OBJ:
			return &Boolean{Value: !obj.(*Boolean).Value}

		default:
			log.Panic("--- illegal unary expression ---")
			log.Printf("%T, %v", node.Right, node.Right)
			log.Printf("unary obj: %v", obj)

		}
	case *InfixExpression:

		left := e.eval(node.Left)
		right := e.eval(node.Right)
		// log.Println("---- infix expression ----")
		// log.Printf("%T, %T", node.Left, node.Right)
		// log.Printf("%v, %v", left, right)
		if left.Type() == right.Type() {
			switch left.Type() {
			case INTEGER_OBJ:
				l := left.(*Integer).Value
				r := right.(*Integer).Value
				return e.evalInfixExpressionInteger(node.Operator, l, r)
			case STRING_OBJ:
				return e.evalInfixExpressionString(node.Operator, left.String(), right.String())
			case BOOLEAN_OBJ:
				l := left.(*Boolean).Value
				r := right.(*Boolean).Value
				return e.evalInfixExpressionBoolean(node.Operator, l, r)
			}
		} else {
			log.Panic("different types in an expression")
		}
	default:
		log.Printf("evalStatement ---> illegal type - - %T", node)
		result = e.eval(node)
		// case *TernaryStatement:
		// 	log.Println("evaluater ---> Ternary")
		// 	result = e.evalTernary(node)
	}

	return result
}
func (e *Evaluater) evalDefine(node *DefineExpression) Object {
	obj := e.eval(node.Expr)
	e.env.Set(node.Ident.Key, obj)
	log.Printf("after Define \n")
	for k, v := range e.env.store {
		log.Println(k, v)
	}
	log.Printf("store end \n")
	return obj
}
func (e *Evaluater) evalAssign(node *AssignExpression) Object {
	if _, err := e.env.Get(node.Ident.Key); err != nil {
		log.Panic(err)
	}
	obj := e.eval(node.Expr)
	e.env.Set(node.Ident.Key, obj)
	return obj
}

func (e *Evaluater) evalArray(node *ArrayLiteral) Object {
	exprs := e.evalExprs(node.Elements)
	array := &Array{Elements: exprs}
	return array
}

func (e *Evaluater) evalIndex(left, idx Object) Object {
	var result Object
	log.Printf("evalIndex : (%v , %v)", left, idx)
	switch {
	case left.Type() == ARRAY_OBJ && idx.Type() == INTEGER_OBJ:
		elmts := left.(*Array).Elements
		i := idx.(*Integer).Value
		if l := len(elmts); i >= 0 && i < l {
			result = elmts[i]
		} else {
			log.Panicf("array index out of range! expect [%d, %d), got %d\n", 0, l, i)
		}
	}
	return result
}
func (e *Evaluater) evalIf(node *IfExpression) Object {
	for i, cnd := range node.conditions {
		if e.isTure(cnd) {
			return e.evalBlock(node.executes[i])
		}
	}
	return nil
}

func (e *Evaluater) evalWhile(node *WhileExpression) Object {
	var result Object
	for e.isTure(node.Condition) {
		result = e.evalBlock(node.Execute)
	}
	return result
}
func (e *Evaluater) evalTernary(node *TernaryExpression) Object {
	if e.isTure(node.condition) {
		return e.eval(node.left)
	}
	return e.eval(node.right)
}

func (e *Evaluater) evalBlock(node *BlockExpression) Object {
	var result Object
	for _, stmt := range node.Statements {
		result = e.eval(stmt)
	}
	return result
}

func (e *Evaluater) evalInfixExpressionInteger(op string, l, r int) Object {
	switch op {
	case "+":
		return &Integer{Value: l + r}
	case "-":
		return &Integer{Value: l - r}
	case "/":
		return &Integer{Value: l / r}
	case "*":
		return &Integer{Value: l * r}
	case "%":
		return &Integer{Value: l % r}
	case ">":
		return &Boolean{Value: l > r}
	case "<":
		return &Boolean{Value: l < r}
	case ">=":
		return &Boolean{Value: l >= r}
	case "<=":
		return &Boolean{Value: l <= r}
	case "==":
		return &Boolean{Value: l == r}
	case "!=":
		return &Boolean{Value: l != r}
	default:
		log.Panic("illegal operator for integer")
	}
	return nil
}

func (e *Evaluater) evalInfixExpressionBoolean(op string, l, r bool) Object {
	switch op {
	case "==":
		return &Boolean{Value: l == r}
	case "!=":
		return &Boolean{Value: l != r}
	default:
		log.Panic("illegal operator for integer")
	}
	return nil
}

func (e *Evaluater) evalInfixExpressionString(op, l, r string) Object {
	switch op {
	case "+":
		return &String{Value: l + r}
	default:
		log.Panic("illegal operator for string")
	}
	return nil
}

func (e *Evaluater) evalExprs(exprs []Expression) []Object {
	objs := []Object{}
	log.Println("---- eval Exprs - - -")
	log.Println(exprs)
	for _, expr := range exprs {
		objs = append(objs, e.eval(expr))
	}
	log.Println(objs)

	return objs
}

func (e *Evaluater) applyFunction(fn Object, args []Object) Object {
	var obj Object
	log.Println("--------- apply function --------")
	log.Println("fn ", fn, " - ", "args ", args)
	switch fn := fn.(type) {
	case *Function:
		e.env = NewFunctionEnvirontment(e.env, fn, args)
		obj = e.evalBlock(fn.Body)
		e.env = e.env.outter
	case *Builtin:
		log.Printf("%T %v", fn, fn)

		obj = fn.Fn(args...)

	}
	return obj
}

// ==================== helper functions
func (e *Evaluater) isTure(node ASTNode) bool {
	res := e.eval(node).(*Boolean)
	return res.Value == true
}
