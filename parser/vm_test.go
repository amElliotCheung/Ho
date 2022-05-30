package interpreter

import (
	"io"
	"log"
	"strings"
	"testing"
)

func TestVM_Run1(t *testing.T) {
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
	a := 1
	b := 2
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
	a := 100
	b := -100
	if a > b {
		a
	} else if a < b {
		b
	} else {
		0
	}
	a = -1000
	if a > b {
		100
	} else if a < b {
		-100
	} else {
		0
	}
	b = -1000
	if a > b {
		100
	} else if a < b {
		-100
	} else {
		0
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
}
func TestVM_Run5(t *testing.T) {
	input := `
	n := 1
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
	n := 1
	sum := 0
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

func TestVM_Run7(t *testing.T) {
	input := `
	n := 2
	add := func (x, y) {
		x+y
	}
	add(1,2)
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

	compiler.show()

	vm := NewVM(compiler.bytecode())
	vm.Run()

}

func TestVM_Run8(t *testing.T) {
	log.SetOutput(io.Discard)
	input := `
	add := func(x,y) {
		x
	} hope {
		-100, 100 -> 0
		0,0 -> 0
		0, 1 -> 1
	}
	add(1,0)
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

	compiler.show()

	vm := NewVM(compiler.bytecode())
	vm.Run()

}
