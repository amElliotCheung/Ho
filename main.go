package main

import (
	"log"
	"os"
	"stone/parser"
)

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||\\n|[[:punct:]])?`

func main() {
	file, err := os.Open("./sourcecode.txt")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	lexer := parser.NewLexer(file, regexPat)
	for tk := lexer.Read(); tk != parser.EOF; tk = lexer.Read() {
		log.Println(tk)
	}
	lexer.ShowQueue()
}
