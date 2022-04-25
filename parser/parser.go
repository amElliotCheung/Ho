package interpreter

import (
	"fmt"
	"log"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[string]int{
	"==": EQUALS,
	"!=": EQUALS,
	"<":  LESSGREATER,
	">":  LESSGREATER,
	"+":  SUM,
	"-":  SUM,
	"/":  PRODUCT,
	"%":  PRODUCT,
	"*":  PRODUCT,
	"(":  CALL,
}

type (
	prefixParseFn func() (ASTNode, error)
	infixParseFn  func(ASTNode) (ASTNode, error)
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

	p.prefixParser[IDENT] = p.parseIdentifier
	p.prefixParser[INT] = p.parseInteger
	p.prefixParser[LPAREN] = p.parseGroupedExpression
	for _, op := range []string{PLUS, MINUS, SLASH, ASTERISK, LT, GT, EQ, NEQ, MOD} {
		p.infixParser[op] = p.parseInfixExpression
	}

	p.advance()
	return p
}

// ================== parse functions
func (p *Parser) Parse(res chan ASTNode) (ASTNode, error) {
	prog := &ASTList{}
	for p.cur != EOF {
		log.Println("parse : cur and next are ", p.cur, p.next)
		if p.cur == EOL {
			p.advance()
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		// debug
		log.Println("parse : ", p.lexer.lineNo, "th line : ", stmt.String())
		// advance
		p.advance()
		prog.addChild(stmt)
		res <- stmt
	}
	close(res)
	return prog, nil
}
func (p *Parser) parseStatement() (ASTNode, error) {
	var stmt ASTNode
	var err error
	switch p.cur.(type) {
	case IdToken:
		if p.checkNext("=") {
			log.Println("id=, enter Assign")
			stmt, err = p.parseAssign()
		} else if p.checkCur("while") {
			log.Println("while, enter while")
			stmt, err = p.parseWhile()
		} else if p.checkCur("if") {
			log.Println("if, enter if")
			stmt, err = p.parseIf()
		} else {
			log.Println("id, enter Expression")
			stmt, err = p.parseExpression(LOWEST)
		}
	case NumToken:
		log.Println("number, enter Expression")
		stmt, err = p.parseExpression(LOWEST)
	case OpToken:
		log.Println("operator, enter Expression")
		stmt, err = p.parseExpression(LOWEST)
	default:
		log.Println("no match for ", p.cur.GetType(), p.cur)
	}

	if err != nil {
		return nil, err
	}
	return stmt, nil
}
func (p *Parser) parseAssign() (ASTNode, error) {
	assign := ASTList{Token: p.next}        // =
	assign.addChild(&ASTLeaf{Token: p.cur}) // identifier
	p.advance()
	p.skip("=")

	expr, err := p.parseExpression(LOWEST) //
	if err != nil {
		return nil, err
	}
	assign.addChild(expr)
	return &AssignStatement{ASTList: assign}, nil
}
func (p *Parser) parseExpression(precedence int) (ASTNode, error) {
	parser, ok := p.prefixParser[p.cur.GetType()]
	if !ok {
		return nil, fmt.Errorf("no prefix function for %v", p.cur.GetType())
	}
	left, _ := parser() // ASTNode, error
	for p.next != EOL && precedence < p.peekPrecedence() {
		infix, ok := p.infixParser[p.next.GetType()]
		if !ok {
			return left, fmt.Errorf("no infix function for %v", p.next.GetType())
		}
		p.advance()
		left, _ = infix(left)
	}
	return left, nil
}
func (p *Parser) parseGroupedExpression() (ASTNode, error) {
	p.skip("(")
	expr, err := p.parseExpression(LOWEST)
	if p.checkNext(")") {
		p.advance()
		return expr, err
	} else {
		return nil, fmt.Errorf("( ) don't match")
	}
}
func (p *Parser) parseInfixExpression(left ASTNode) (ASTNode, error) {
	expr := ASTList{Token: p.cur}
	expr.addChild(left)

	precedence := p.curPrecedence()
	p.advance()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	expr.addChild(right)

	return &BinaryExpr{ASTList: expr}, nil
}

func (p *Parser) parseIf() (ASTNode, error) {
	node := ASTList{Token: p.cur}
	p.advance()
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	p.advance()
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	node.addChild(expr)
	node.addChild(block)
	// ------------------- else
	for p.checkCur("else") {
		if p.checkNext("if") {
			p.advance()
			p.advance()
			expr, err = p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			p.advance()
			block, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			node.addChild(expr)
			node.addChild(block)
		} else {
			p.advance()
			block, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			node.addChild(&ASTLeaf{Token: NewIdToken(p.lexer.lineNo, "true")}) // true condition
			node.addChild(block)
		}
	}
	return &IfStatement{ASTList: node}, nil
}

func (p *Parser) parseWhile() (ASTNode, error) {
	node := ASTList{Token: p.cur}
	p.advance()
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	p.advance()
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	node.addChild(expr)
	node.addChild(block)

	return &WhileStatement{ASTList: node}, nil
}

func (p *Parser) parseBlock() (ASTNode, error) {
	p.skip("{")
	block := &ASTList{}
	for p.cur.GetText() != "}" {
		log.Println("cur and next are ", p.cur, p.next)

		if p.cur == EOL {
			p.advance()
			continue
		}

		stmt, err := p.parseStatement()

		if err != nil {
			return nil, err
		}
		// debug
		log.Println(p.lexer.lineNo, "th line : ", stmt.String())
		// advance
		p.advance()
		block.addChild(stmt)

	}
	if p.checkCur("}") {
		p.advance()
	}
	return block, nil
}

func (p *Parser) parseIdentifier() (ASTNode, error) {
	return &ASTLeaf{Token: p.cur}, nil
}

func (p *Parser) parseInteger() (ASTNode, error) {
	return &ASTLeaf{Token: p.cur}, nil
}
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.next.GetType()]; ok {
		return p
	}

	return LOWEST
}

// ======================= helper functions

func (p *Parser) skip(s string) {
	if p.checkCur(s) {
		p.advance()
	} else {
		log.Fatalln(p.lexer.lineNo, "expect ", s, "got", p.cur.GetText())
	}
}

func (p *Parser) advance() {
	p.cur = p.lexer.Read()
	p.next = p.lexer.Peek(0)
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.cur.GetType()]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) checkCur(expt string) bool {
	return p.cur.GetText() == expt
}

func (p *Parser) checkNext(expt string) bool {
	return p.next.GetText() == expt
}
