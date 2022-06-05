package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"stone/interpreter"
)

// the package name must be the same as foldername
func main() {
	log.SetOutput(io.Discard)

	filename := flag.String("f", "sourcecode.txt", "the file containing source code")
	productive := flag.Bool("p", false, "the compiler would ignore hope block if this variable is true")
	flag.Parse()

	input, err := os.ReadFile(*filename)
	if err != nil {
		fmt.Println(err)
	}

	in := strings.NewReader(string(input))
	lexer := interpreter.NewLexer(in)
	parser := interpreter.NewParser(lexer)
	node, err := parser.Parse(nil)
	if err != nil {
		log.Println(err)
	}
	compiler := interpreter.NewCompiler(*productive)
	compiler.Compile(node)

	vm := interpreter.NewVM(compiler.Bytecode())
	vm.Run()

}

// //================== test lexer
// func lexer_test(filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer file.Close()

// 	lexer := interpreter.NewLexer(file)
// 	for tk := lexer.Read(); tk != interpreter.EOF; tk = lexer.Read() {
// 		log.Println(tk)
// 	}
// }

// func parser_test(filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer file.Close()

// 	lexer := interpreter.NewLexer(file, regexPat)
// 	parser := interpreter.NewParser(lexer)

// 	root, err := parser.Parse(nil)
// 	if err != nil {
// 		log.Println("error:\t", err)
// 	} else {
// 		log.Println(root)
// 	}

// }

// func evaluater_test(filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer file.Close()

// 	lexer := interpreter.NewLexer(file, regexPat)
// 	parser := interpreter.NewParser(lexer)
// 	c := make(chan interpreter.ASTNode)
// 	go parser.Parse(c)
// 	evaluater := interpreter.NewEvaluater(c)
// 	result := evaluater.Eval()
// 	log.Println("=============  evaluater final result  ============")
// 	log.Println("-- ", result.Type(), result.Inspect(), "--")
// }
