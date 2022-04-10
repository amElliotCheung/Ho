package lexer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"stone/token"
	"strconv"
)

type Lexer struct {
	pat     *regexp.Regexp // regular expression
	scanner *bufio.Scanner
	queue   []token.Token // list of tokens
	lineNo  int
	hasMore bool
}

func NewLexer(in io.Reader, regexpPat string) *Lexer {
	return &Lexer{
		pat:     regexp.MustCompile(regexpPat),
		scanner: bufio.NewScanner(in),
		queue:   make([]token.Token, 2),
		lineNo:  0,
		hasMore: true,
	}
}

func (l *Lexer) Readline() {
	// ------ read a line
	line := ""
	if l.scanner.Scan() {
		line = l.scanner.Text()
		l.lineNo += 1
		log.Printf("%v--%v", line, line[len(line)-1]) // output
	} else {
		l.hasMore = false
		return
	}
	// ------ match a token
	low := 0
	for low < len(line) {
		if s := l.pat.FindString(line[low:]); s != "" { // a token is matched
			low += len(s)
			l.AddToken(s)
		} else {
			log.Fatalln("bad token at line ", l.lineNo)
		}
	}
	l.AddToken("\\n")
}

func (l *Lexer) AddToken(str string) {
	matches := l.pat.FindAllStringSubmatch(str, -1)[0]
	m1 := matches[1] // the first par match
	fmt.Println("the match:\n\t", m1)
	if m1 == "" || m1 == `\n` || matches[2] != "" { // empty or \n or comment
		return
	}
	var tk token.Token
	if matches[3] != "" { // number
		val, _ := strconv.Atoi(matches[3])
		tk = token.NewNumToken(l.lineNo, val)
	} else if matches[4] != "" { // string
		tk = token.NewStrToken(l.lineNo, l.toStringLiteral(matches[4]))
	} else { // identifier
		tk = token.NewIdToken(l.lineNo, m1)
	}
	l.queue = append(l.queue, tk)
}

func (l *Lexer) toStringLiteral(s string) string {
	return fmt.Sprint("%s", s)
}

// a helper function to debug
func (l *Lexer) ShowQueue() {
	for _, item := range l.queue {
		if item != nil {
			switch v := item.(type) {
			case token.NumToken:
				fmt.Printf("Number: %v\n", v.GetText())
			case token.StrToken:
				fmt.Printf("String: %v\n", v.GetText())
			case token.IdToken:
				fmt.Printf("Ident: %v\n", v.GetText())
			default:
				// And here I'm feeling dumb. ;)
				fmt.Printf("I don't know, ask stackoverflow.")
			}

		}
	}
}
