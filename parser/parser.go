package interpreter

import (
	"fmt"
	"log"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	ASSIGNPRE
	QUESTIONMARK // ? :
	EQUALS       // ==
	LESSGREATER  // > or <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X or !X
	CALL         // myFunction(X)
	INDEX        // array[index]
)

var precedences = map[string]int{
	"==": EQUALS,
	"!=": EQUALS,
	"<":  LESSGREATER,
	">":  LESSGREATER,
	"<=": LESSGREATER,
	">=": LESSGREATER,
	"+":  SUM,
	"-":  SUM,
	"/":  PRODUCT,
	"%":  PRODUCT,
	"*":  PRODUCT,
	"(":  CALL,
	"?":  QUESTIONMARK,
	"=":  ASSIGNPRE,

	"[": INDEX,
}

type (
	prefixParseFn func() (Expression, error)
	infixParseFn  func(Expression) (Expression, error)
)

type Parser struct {
	lexer        *Lexer
	cur          Token
	next         Token
	prefixParser map[string]prefixParseFn
	infixParser  map[string]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer:        l,
		prefixParser: make(map[string]prefixParseFn),
		infixParser:  make(map[string]infixParseFn),
	}

	p.prefixParser[IDENTIFIER] = p.parseIdentifier
	p.prefixParser[INTEGER] = p.parseInteger
	p.prefixParser[BOOLEAN] = p.parseBoolean
	p.prefixParser[STRING] = p.parseString
	p.prefixParser[LBRACKET] = p.parseArray
	p.prefixParser[LPAREN] = p.parseGroupedExpression
	p.prefixParser[BANG] = p.parseUnaryExpression
	p.prefixParser[MINUS] = p.parseUnaryExpression
	p.prefixParser[FUNCTION] = p.parseFunction

	for _, op := range []string{PLUS, MINUS, SLASH, ASTERISK, LT, GT, LTE, GTE, EQ, NEQ, MOD, ASSIGN} {
		p.infixParser[op] = p.parseInfixExpression
	}

	p.infixParser["?"] = p.parseTernaryExpression
	p.infixParser[LPAREN] = p.parseCallExpression
	p.infixParser[LBRACKET] = p.parseIndexExpression
	p.advance()
	return p
}

// ================== parse functions
func (p *Parser) Parse(res chan Statement) (ASTNode, error) {
	prog := &Program{}
	for p.cur != EOF {
		// // log.Println("parse : cur and next are ", p.cur, p.next)
		if p.cur == EOL {
			p.advance()
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		// debug
		// log.Println("parse : ", p.lexer.lineNo, "th line : ", stmt.String())
		// advance
		p.advance()
		prog.Statements = append(prog.Statements, stmt)
		if res != nil {
			res <- stmt
		}

	}
	if res != nil {
		close(res)
	}

	return prog, nil
}
func (p *Parser) parseStatement() (Statement, error) {
	var stmt Statement
	var err error

	switch p.cur.Literal() {
	default:
		stmt, err = p.parseExpressionStatement()
	}

	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseExpressionStatement() (Expression, error) {
	switch p.cur.Literal() {
	case "while":
		return p.parseWhileExpression()
	case "if":
		return p.parseIfExpression()
	default:
		switch p.next.Literal() {
		case ":=":
			return p.parseDefineExpression()
		case "=":
			return p.parseAssignExpression()
		default:
			return p.parseExpression(LOWEST)
		}

	}
}
func (p *Parser) parseDefineExpression() (Expression, error) {
	ds := &DefineExpression{}                           // =
	ds.Ident = &IdentifierLiteral{Key: p.cur.Literal()} // identifier
	p.advance()
	p.skip(":=")

	ds.Expr, _ = p.parseExpression(LOWEST) //
	return ds, nil
}
func (p *Parser) parseAssignExpression() (Expression, error) {
	assign := &AssignExpression{}                           // =
	assign.Ident = &IdentifierLiteral{Key: p.cur.Literal()} // identifier
	p.advance()
	p.skip("=")

	assign.Expr, _ = p.parseExpression(LOWEST) //
	return assign, nil
}

func (p *Parser) parseIfExpression() (Expression, error) {
	ie := &IfExpression{conditions: make([]Expression, 0),
		executes: make([]*BlockExpression, 0)}

	p.skip("if")
	log.Println("----- parseIfExpression -----  ", p.cur.Type(), p.cur.Literal())
	cnd, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	p.advance()
	block, err := p.parseBlockExpression()
	if err != nil {
		return nil, err
	}
	ie.addPair(cnd, block)
	// ------------------- else
	for p.checkCur("else") {
		p.advance()
		if p.checkCur("if") {
			p.advance()
			cnd, err = p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			p.advance()
		} else {
			cnd = &BooleanLiteral{Key: true}
		}

		block, err := p.parseBlockExpression()
		if err != nil {
			return nil, err
		}
		ie.addPair(cnd, block)
	}
	return ie, nil
}

func (p *Parser) parseWhileExpression() (ASTNode, error) {
	we := &WhileExpression{}
	p.advance()
	we.Condition, _ = p.parseExpression(LOWEST)

	p.advance()
	we.Execute, _ = p.parseBlockExpression()

	return we, nil
}

func (p *Parser) parseBlockExpression() (*BlockExpression, error) {
	// log.Println("block!\t")
	p.skip("{")
	block := &BlockExpression{Statements: make([]Statement, 0)}
	for !p.checkCur(RBRACE) {
		// log.Println("cur and next are ", p.cur, p.next)

		if p.cur == EOL {
			p.advance()
			continue
		}

		stmt, err := p.parseStatement()

		if err != nil {
			return nil, err
		}
		// debug
		// log.Println(p.lexer.lineNo, "th line : ", stmt.String())
		// advance
		p.advance()
		block.Statements = append(block.Statements, stmt)

	}
	p.skip("}")
	return block, nil
}
func (p *Parser) parseExpression(precedence int) (Expression, error) {

	tp := p.cur.Type()

	if tp == OPERATOR {
		tp = p.cur.Literal()
	}
	parser, ok := p.prefixParser[tp]
	if !ok {
		return nil, fmt.Errorf("no prefix function for %+v", p.cur.Literal())
	}
	left, _ := parser() // Expression, error
	for p.next.Literal() != "EOL" && precedence < p.peekPrecedence() {
		log.Println("---- parseExpression ----", p.cur, p.next)

		tp = p.next.Type()
		if tp == OPERATOR {
			tp = p.next.Literal()
		}
		infix, ok := p.infixParser[tp]
		if !ok {
			return left, fmt.Errorf("no infix function for %v", p.next.Literal())
		}
		p.advance()
		left, _ = infix(left)

	}
	log.Println("---- parseExpression end----", p.cur, p.next)

	return left, nil
}

func (p *Parser) parseUnaryExpression() (Expression, error) {
	log.Printf("========= unary expression========%T %v", p.cur, p.cur)
	ue := &UnaryExpression{
		Operator: p.cur.Literal(),
	}
	p.advance()
	ue.Right, _ = p.parseExpression(PREFIX)
	return ue, nil
}
func (p *Parser) parseGroupedExpression() (Expression, error) {
	p.skip("(")
	expr, _ := p.parseExpression(LOWEST)
	p.advance()
	// (2+4/2)*(2+3)
	// p.skip(")")
	return expr, nil
}
func (p *Parser) parseInfixExpression(left Expression) (Expression, error) {
	expr := &InfixExpression{
		Operator: p.cur.Literal(),
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.advance()
	expr.Right, _ = p.parseExpression(precedence)
	return expr, nil
}
func (p *Parser) parseTernaryExpression(condition Expression) (Expression, error) {
	// log.Println("Now I am in Ternary?")
	p.skip("?")
	ternary := &TernaryExpression{
		condition: condition,
	}
	ternary.left, _ = p.parseExpression(LOWEST)
	p.advance()
	p.skip(":")
	ternary.right, _ = p.parseExpression(LOWEST)
	return ternary, nil
}

func (p *Parser) parseHopeBlock() (*HopeBlock, error) {
	p.skip("{")
	hopeBlock := &HopeBlock{
		HopeExpressions: make([]HopeExpression, 0),
	}
	for !p.checkCur(RBRACE) {
		if p.cur == EOL {
			p.advance()
			continue
		}
		hpe := HopeExpression{}
		hpe.Parameters, _ = p.parseExpressionList("->")
		log.Printf("++\n\nafter parse parameters %v %v\n\n++", p.cur, p.next)
		p.skip("->")
		// should advance() after parseExpression
		hpe.Expected, _ = p.parseExpression(LOWEST)
		p.advance()
		log.Printf("++\n\nafter parse answer %v %v\n\n++", p.cur, p.next)
		hopeBlock.HopeExpressions = append(hopeBlock.HopeExpressions, hpe)
	}
	p.skip("}")
	return hopeBlock, nil
}
func (p *Parser) parseIndexExpression(left Expression) (Expression, error) {
	expr := &IndexExpression{
		Left: left,
	}
	p.skip(LBRACKET)
	expr.Index, _ = p.parseExpression(LOWEST)
	p.advance()
	// p.skip(RBRACKET)
	return expr, nil
}
func (p *Parser) parseCallExpression(left Expression) (Expression, error) {
	ce := &CallExpression{
		Function: left,
	}
	p.skip(LPAREN)
	log.Println("-- CallExpression --", p.cur, p.next)

	ce.Arguments, _ = p.parseExpressionList(RPAREN)
	// delete this so that len(a) <= 10 could work
	// p.skip(RPAREN)
	return ce, nil
}

func (p *Parser) parseExpressionList(end string) ([]Expression, error) {
	list := []Expression{}
	if p.checkCur(end) { // no identifier
		return list, nil
	}
	for {
		log.Printf("before parse expressionList p.cur=%v, p.next=%v\n", p.cur, p.next)
		expr, _ := p.parseExpression(LOWEST)
		log.Printf("after parse expressionList p.cur=%v, p.next=%v\n", p.cur, p.next)

		list = append(list, expr)
		p.advance()
		log.Printf("expressionList: after advance p.cur=%v, p.next=%v\n", p.cur, p.next)

		if p.checkCur(end) {
			return list, nil
		} else if p.checkCur(COMMA) {
			p.advance()
		} else {
			return nil, fmt.Errorf("function parameters mismatch")
		}
	}
}

// ================= parse leaves

func (p *Parser) parseArray() (Expression, error) {
	array := &ArrayLiteral{Elements: make([]Expression, 0)}
	p.skip(LBRACKET)
	array.Elements, _ = p.parseExpressionList(RBRACKET)
	p.skip(RBRACKET)
	return array, nil
}

func (p *Parser) parseFunction() (Expression, error) {
	p.skip("func")
	paras, _ := p.parseIdentifierList()
	p.skip(RPAREN)
	exec, _ := p.parseBlockExpression()

	var hopes *HopeBlock
	hopes = nil
	if p.checkCur("hope") {
		p.skip("hope")
		hopes, _ = p.parseHopeBlock()
	}
	return &FunctionLiteral{
		Parameters: paras,
		Execute:    exec,
		Hopes:      hopes,
	}, nil
}
func (p *Parser) parseIdentifierList() ([]*IdentifierLiteral, error) {
	p.skip(LPAREN)
	list := []*IdentifierLiteral{}
	if p.checkCur(RPAREN) { // no identifier
		return list, nil
	}
	for {
		ident := &IdentifierLiteral{Key: p.cur.Literal()}
		list = append(list, ident)
		p.advance()
		if p.checkCur(RPAREN) {
			return list, nil
		} else if p.checkCur(COMMA) {
			p.advance()
		}
	}
}

// ========== parse leaves
func (p *Parser) parseString() (Expression, error) {
	return &StringLiteral{Key: p.cur.Literal()}, nil
}

func (p *Parser) parseIdentifier() (Expression, error) {
	return &IdentifierLiteral{Key: p.cur.Literal()}, nil
}

func (p *Parser) parseInteger() (Expression, error) {
	val, err := strconv.Atoi(p.cur.Literal())
	if err != nil {
		return nil, err
	}
	return &IntegerLiteral{Key: val}, nil
}

func (p *Parser) parseBoolean() (Expression, error) {
	key := false
	if p.cur.Literal() == "true" {
		key = true
	}
	return &BooleanLiteral{Key: key}, nil
}

//  ============ helper functions

func (p *Parser) peekPrecedence() int {
	t := p.next.Type()
	if t == OPERATOR {
		t = p.next.Literal()
	}
	if p, ok := precedences[t]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) skip(s string) {
	if p.checkCur(s) {
		p.advance()
	} else {
		log.Panicf("skip: want %s, got <%s %s>", s, p.cur.Type(), p.cur.Literal())

	}
}

func (p *Parser) advance() {
	p.cur = p.lexer.Read()
	p.next = p.lexer.Peek(0)
}

func (p *Parser) curPrecedence() int {
	t := p.cur.Type()
	if t == OPERATOR {
		t = p.cur.Literal()
	}
	if p, ok := precedences[t]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) checkCur(expt string) bool {
	log.Printf("checkCur %v, %v\n", p.cur.Literal(), expt)
	return p.cur.Literal() == expt
}

func (p *Parser) checkNext(expt string) bool {
	return p.next.Literal() == expt
}
