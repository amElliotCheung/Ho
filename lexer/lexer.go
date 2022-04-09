package lexer

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"stone/token"
)

type Lexer struct {
	pat     *regexp.Regexp // regular expression
	scanner *bufio.Scanner
	queue   []token.Token // list of tokens
	lineNo  int
}

func NewLexer(in io.Reader, regexpPat string) *Lexer {
	return &Lexer{
		pat:     regexp.MustCompile(regexpPat),
		scanner: bufio.NewScanner(in),
		queue:   make([]token.Token, 2),
		lineNo:  0,
	}
}

func (l Lexer) readline() {
	if l.scanner.Scan() {
		line := l.scanner.Text()
		log.Println(line)
	}
}
