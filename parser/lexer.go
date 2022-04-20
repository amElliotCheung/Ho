package parser

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
)

type Lexer struct {
	pat     *regexp.Regexp // regular expression
	scanner *bufio.Scanner
	queue   []Token // list of tokens
	lineNo  int
	hasMore bool
}

func NewLexer(in io.Reader, regexpPat string) *Lexer {
	return &Lexer{
		pat:     regexp.MustCompile(regexpPat),
		scanner: bufio.NewScanner(in),
		queue:   make([]Token, 0),
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
		//log.Printf("%v--%v", line, line[len(line)-1]) // output
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
	l.queue = append(l.queue, EOL)
}

func (l *Lexer) fillQueue(i int) bool {
	for i >= len(l.queue) {
		if l.hasMore == false {
			return false
		}
		l.Readline()
	}
	return true
}

func (l *Lexer) Peek(i int) Token {
	if l.fillQueue(i) == false {
		return EOF
	}
	return l.queue[i]
}

func (l *Lexer) Read() Token {
	if l.fillQueue(0) {
		tk := l.queue[0]
		l.queue = l.queue[1:]
		return tk
	}
	return EOF

}

func (l *Lexer) AddToken(str string) {
	matches := l.pat.FindAllStringSubmatch(str, -1)[0]
	m1 := matches[1] // the first par match

	// for i, m := range matches {
	// 	log.Println("matches : ", i, m)
	// }
	if m1 == "" || matches[2] != "" { // empty or \n or comment
		return
	}
	var tk Token
	if matches[3] != "" { // number
		val, _ := strconv.Atoi(matches[3])
		tk = NewNumToken(l.lineNo, val)
	} else if matches[4] != "" { // string
		tk = NewStrToken(l.lineNo, l.toStringLiteral(matches[4]))
	} else if matches[6] != "" { // identifier
		tk = NewIdToken(l.lineNo, matches[6])
	} else {
		tk = NewOpToken(l.lineNo, m1)
	}
	l.queue = append(l.queue, tk)
	// log.Printf("add %T type token", tk)
}

func (l *Lexer) toStringLiteral(s string) string {
	return fmt.Sprint("%s", s)
}

// a helper function to debug
func (l *Lexer) ShowQueue() {
	for _, item := range l.queue {
		if item != nil {
			switch v := item.(type) {
			case NumToken:
				fmt.Printf("Number: %v\n", v.GetText())
			case StrToken:
				fmt.Printf("String: %v\n", v.GetText())
			case IdToken:
				fmt.Printf("Ident: %v\n", v.GetText())
			default:
				// And here I'm feeling dumb. ;)
				fmt.Printf("I don't know, ask stackoverflow.")
			}

		}
	}
}