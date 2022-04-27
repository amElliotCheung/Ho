package main

import (
	"log"
	"os"
	interpreter "stone/parser"
)

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|([A-Za-z]\w*)|(\+|-|\*|/|%|==|!=|<|<=|>|>=|&&|\|\||\\n|\?|:|[[:punct:]]))?`

func main() {
	log.SetFlags(0) // no prefix
	// lexer_test(filename)
	// lexer_test("./sourcecode.txt")
	// parser_test("./sourcecode.txt")
	// go evaluater_test("./evaluation.txt")
	evaluater_test("./fact.txt")

}

//================== test lexer
func lexer_test(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	lexer := interpreter.NewLexer(file, regexPat)
	for tk := lexer.Read(); tk != interpreter.EOF; tk = lexer.Read() {
		log.Println(tk)
	}
}

func parser_test(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	lexer := interpreter.NewLexer(file, regexPat)
	parser := interpreter.NewParser(lexer)

	root, err := parser.Parse(nil)
	if err != nil {
		log.Println("error:\t", err)
	} else {
		log.Println(root)
	}

}

func evaluater_test(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	lexer := interpreter.NewLexer(file, regexPat)
	parser := interpreter.NewParser(lexer)
	c := make(chan interpreter.ASTNode)
	go parser.Parse(c)
	evaluater := interpreter.NewEvaluater(c)
	result := evaluater.Eval()
	log.Println("=============  evaluater final result  ============")
	log.Println("-- ", result.Type(), result.Inspect(), "--")
}
