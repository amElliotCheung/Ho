package interpreter

import (
	"log"
	"strings"
	"testing"
)

func TestVM_Run(t *testing.T) {
	input := `
	1+2
	2+3+4
	5*2%7
	(2+4/2)*(2+3)
	!true	
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler()
	compiler.Compile(node)
	vm := NewVM(compiler.bytecode())
	vm.Run()
}
func TestVM_Run2(t *testing.T) {
	input := `
	("one" + "two" == "onetwo")	
	!("one" + "two" == "onetwo")	
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler()
	compiler.Compile(node)
	vm := NewVM(compiler.bytecode())
	vm.Run()
}

func TestVM_Run3(t *testing.T) {
	input := `
	a = 1
	b = 2
	a + b
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler()
	compiler.Compile(node)
	vm := NewVM(compiler.bytecode())
	vm.Run()
}
func TestVM_Run4(t *testing.T) {
	input := `
	a = 3
	b = 1
	if a > b {
		2
	}
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler()
	compiler.Compile(node)
	vm := NewVM(compiler.bytecode())
	vm.Run()
	// compiler.show()
}
func TestVM_Run5(t *testing.T) {
	input := `
	n = 1
	while n < 10 {
		n = n + 1
	}
	n
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler()
	compiler.Compile(node)
	vm := NewVM(compiler.bytecode())
	vm.Run()

	// log.Println("======   compiler instruction   ======")
	// compiler.show()
}
func TestVM_Run6(t *testing.T) {
	input := `
	n = 1
	sum = 0
	while n <= 100 {
		sum = sum + n
		n = n + 1
	}
	sum
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in, regexPat)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler()
	compiler.Compile(node)
	vm := NewVM(compiler.bytecode())
	vm.Run()

	// log.Println("======   compiler instruction   ======")
	// compiler.show()
}
