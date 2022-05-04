package interpreter

import (
	"strings"
	"testing"
)

func TestEvaluater_Eval(t *testing.T) {
	input := `adder := func (x, y) {
		x+y
	}
	a := -1
	b := 1
	adder(a, b)`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	c := make(chan Statement)
	go parser.Parse(c)
	e := NewEvaluater(c)
	e.Eval()

}
