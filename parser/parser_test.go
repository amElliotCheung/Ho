package interpreter

import (
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	input := `a := 0
	a= -1-1-1-1-1*3-5+(-9*3)
	p := 1
	b := "123"
	c := "asdf"
	d := "a\n\n\n"
	afunc := func(){}
	array := [1,b,c,d,afun]
	array[2]
	if p>a{
		array[0]
	} else if p == a {
		array[1]
	} else {
		array[2]
	}`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	root, err := parser.Parse(nil)
	t.Log("================ parse root =============\n ", root)
	t.Log("================ parse ERROR =============\n ", err)

}
func TestParser_Parse2(t *testing.T) {
	input := `adder := func (x, y) {
		x+y
	}
	a := -1
	b := 1
	adder(a, b)`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	root, err := parser.Parse(nil)
	t.Log("================ parse root =============\n ", root)
	t.Log("================ parse ERROR =============\n ", err)

}

func TestParser_Parse3(t *testing.T) {
	input := `a = true
	b = !a`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	root, err := parser.Parse(nil)
	t.Log("================ parse root =============\n ", root)
	t.Log("================ parse ERROR =============\n ", err)

}

func TestParser_Parse4(t *testing.T) {
	input := `a := [1, 2]
	l := len(a)
	while len(a) <= 10 {
		a = append(a, 3)
	}
	l`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	root, err := parser.Parse(nil)
	t.Log("================ parse root =============\n ", root)
	t.Log("================ parse ERROR =============\n ", err)

}

func TestParser_Parse5(t *testing.T) {
	input := `fib := func(n) {
			if n <= 2{
				n
			} else {
				fib(n-1) + fib(n-2)
			}
		} hope {
			0 -> 0	
			1 -> 1
			2 -> 2
		}`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	root, err := parser.Parse(nil)
	t.Log("================ parse root =============\n ", root)
	t.Log("================ parse ERROR =============\n ", err)

}
