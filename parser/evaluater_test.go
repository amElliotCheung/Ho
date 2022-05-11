package interpreter

import (
	"strings"
	"testing"
)

func TestEvaluater_Eval(t *testing.T) {
	input := `
	adder := func (x, y) {
		x+y
	}
	a := -1
	b := 1
	adder(a, b)
	c := "one"
	d := " two"
	adder(c, d)`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	c := make(chan Statement)
	go parser.Parse(c)
	e := NewEvaluater(c)
	e.Eval()
}

func TestEvaluater_Eval2(t *testing.T) {
	input := `
	adder := func (x, y) {
		x+y
	}
	n := 1
	a := 0
	b := 1
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	c := make(chan Statement)
	go parser.Parse(c)
	e := NewEvaluater(c)
	e.Eval()
}
func TestEvaluater_Eval3(t *testing.T) {
	input := `a := true
	b := !a
	b = !!!!!b
	a = b == (1 == 1)
	a = b != (1 == 1)
	1 != 2
	1+1 == 2
	3 > 5
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	c := make(chan Statement)
	go parser.Parse(c)
	e := NewEvaluater(c)
	e.Eval()
}

func TestEvaluater_Eval4(t *testing.T) {
	input := `a := [1, 2]
	while len(a) <= 10 {
		a = append(a, a[len(a)-1] + a[len(a)-2])
	}
	a[9]
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	c := make(chan Statement)
	go parser.Parse(c)
	e := NewEvaluater(c)
	e.Eval()
}
func TestEvaluater_Eval5(t *testing.T) {
	input := `
	fib := func (n) {
		a := [1, 2]
		while len(a) <= n {
			a = append(a,  a[len(a)-1] + a[len(a)-2])
		}
		a[n-1]
	}
	add := func(x, y) {
		x+y
	}
	fib(50)
	add("one", "+two")`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	c := make(chan Statement)
	go parser.Parse(c)
	e := NewEvaluater(c)
	e.Eval()
}
