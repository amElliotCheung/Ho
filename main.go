package main

import (
	"log"
	"os"
	"stone/lexer"
)

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||\\n|[[:punct:]])?`

func main() {
	file, err := os.Open("./sourcecode.txt")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	lexer := lexer.NewLexer(file, regexPat)
	for i := 0; i < 10; i++ {
		lexer.Readline()
	}
	lexer.ShowQueue()
}
