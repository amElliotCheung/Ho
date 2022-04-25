package interpreter

// type Evaluater struct {
// 	c chan ASTNode
// }

// func NewEvaluater(c chan ASTNode) *Evaluater {
// 	return &Evaluater{c: c}
// }
// func (e *Evaluater) Eval() ASTNode {
// 	for node := range e.c {
// 		switch node.(type) {
// 		case *ASTLeaf: // just a token
// 			Token.IsIdentifier()
// 		case *BinaryExpr:
// 			e.evalBinaryExpr(node.(*BinaryExpr))
// 		}
// 	}
// }

// func (e *Evaluater) evalBinaryExpr(expr *BinaryExpr) int {
// 	switch expr.operator() {
// 	case "+":
// 		return e.evalBinaryExpr(expr.left()) + e.evalBinaryExpr(expr.right())
// 	case "-":
// 	case "*":
// 	case "/":
// 	case "%":

// 	}
// }
