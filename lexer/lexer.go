package lexer

import (
	"regexp"
	"stone/token"
)

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||[[:punct:]])?`

type Lexer struct {
	pat    *regexp.Regexp
	input  string
	pos    int
	curpos int
	chr    byte
	queue  []token.Token  // list of tokens
	regex  *regexp.Regexp // regular expression
}
func NewLexer(regexpPat string) *Lexer {
	return &Lexer{
		pat: regexp.MustCompile(regexpPat)
	}
	
}
