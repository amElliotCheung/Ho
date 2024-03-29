package interpreter

import (
	"io"
	"log"
	"os"
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
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)
	vm := NewVM(compiler.Bytecode())
	vm.Run()
}
func TestVM_Run2(t *testing.T) {
	input := `
	("one" + "two" == "onetwo")	
	!("one" + "two" == "onetwo")	
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)
	vm := NewVM(compiler.Bytecode())
	vm.Run()
}

func TestVM_Run3(t *testing.T) {
	input := `
	a := 1
	b := 2
	a + b
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)
	vm := NewVM(compiler.Bytecode())
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
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)
	vm := NewVM(compiler.Bytecode())
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
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)
	vm := NewVM(compiler.Bytecode())
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
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)
	vm := NewVM(compiler.Bytecode())
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
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)

	compiler.show()

	vm := NewVM(compiler.Bytecode())
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
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)

	compiler.show()

	vm := NewVM(compiler.Bytecode())
	vm.Run()

}

func TestVM_Run9(t *testing.T) {
	// log.SetOutput(io.Discard)
	log.SetOutput(os.Stdout)
	input := `
	if 2 > 1 {
		2
	} else{
		1
	}
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)

	compiler.show()

	vm := NewVM(compiler.Bytecode())
	vm.Run()
}

func TestVM_Run10(t *testing.T) {
	log.SetOutput(io.Discard)

	input := `
	fib := func(n) {
		if n <= 2{
			n
		} else {
			fib(n-1) + fib(n-2)
		}
	}
	fib(10)
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)

	compiler.show()

	vm := NewVM(compiler.Bytecode())
	vm.Run()
}

func TestVM_Run11(t *testing.T) {
	// log.SetOutput(os.Stdout)
	log.SetOutput(io.Discard)
	input := `
	fib := func(n) {
		if n <= 2{
			n
		} else {
			fib(n-1) + fib(n-2)
		}
	} hope {
		3 -> 3
		4 -> 5
		6 -> 7
	}
	fib(3)
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)

	compiler.show()

	vm := NewVM(compiler.Bytecode())
	vm.Run()
}
func TestVM_Run12(t *testing.T) {
	// log.SetOutput(os.Stdout)
	log.SetOutput(io.Discard)
	input := `
	fib := func(n) {
		if n <= 2{
			n
		} else {
			fib(n-1) + fib(n-2)
		}
	} hope {
		2 -> 2
		35 -> 2
	}
	add := func(x, y) {
		x+y
	} hope {
		2,3 -> 3
	}
	
	fib(35)
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)

	compiler.show()

	vm := NewVM(compiler.Bytecode())
	vm.Run()
}

func TestVM_Run13(t *testing.T) {
	// log.SetOutput(os.Stdout)
	log.SetOutput(io.Discard)
	input := `
	// a wrong function
	// virtual machine would give warning
	// "want 0, got 1 in the 1st test case"
	add := func(x int, y int) {
		x
	} hope {
		1, -1 -> 0
	}
	sum := func(n int) {
		ans := 0
		i := 1
		while i <= n {
			ans = ans + i
			i = i + 1
		}
		ans
	} hope {
		1 -> 1
		10 -> 55
	}
	sum(2)
	sum(5)
	sum(10)
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)

	compiler.show()

	vm := NewVM(compiler.Bytecode())
	vm.Run()
}

func TestVM_Run14(t *testing.T) {
	// log.SetOutput(os.Stdout)
	log.SetOutput(io.Discard)
	input := `
	fib := func(n int) {
		if n <= 3 {
			n / 2
		} else {
			fib(n-1) + fib(n-2)
		}
	} hope {
		1 -> 0
		2 -> 1
		3 -> 1
		10 -> 34
	}
	add := func(x int, y int) {
		x+y
	} hope {
		2,3 -> 1 // a wrong case
		fuzzing 9
	}

	reserve := func(b bool) {
		!b
	} hope {
		true -> false
		false -> true
	}
	
	reserver(true)
	`
	in := strings.NewReader(input)
	lexer := NewLexer(in)
	parser := NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := NewCompiler(false)
	compiler.Compile(node)

	compiler.show()

	vm := NewVM(compiler.Bytecode())
	vm.Run()
}
