package main

import (
	"log"
	"os"
	interpreter "stone/parser"
)

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|([A-Za-z]\w*)|(\+|-|\*|/|%|==|<|<=|>|>=|&&|\|\||\\n|[[:punct:]]))?`

func main() {
	log.SetFlags(0) // no prefix
	filename := "./sourcecode.txt"
	// lexer_text(filename)
	// lexer_text("./sourcecode.txt")
	parser_text(filename)
	// evaluater_text(filename)
}

//================== test lexer
func lexer_text(filename string) {
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

func parser_text(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	lexer := interpreter.NewLexer(file, regexPat)
	parser := interpreter.NewParser(lexer)

	c := make(chan interpreter.ASTNode)
	go parser.Parse(c)
	if err != nil {
		log.Println("error:\t", err)
	}
	for item := range c {
		log.Println("receive item from channel : ", item)
	}

}

// func evaluater_text(filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	defer file.Close()

// 	lexer := interpreter.NewLexer(file, regexPat)
// 	parser := interpreter.NewParser(lexer)
// 	c := make(chan interpreter.ASTNode)
// 	parser.Parse(c)
// 	evaluater := interpreter.NewEvaluater(c)
// 	evaluater.Eval()
// }
