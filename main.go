package main

import (
	"os"
	"stone/token"
)

const regexPat = `\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||[[:punct:]])?`

func main() {
	file, _ := os.Open("./sourcecode.txt")
	defer file.Close()
	lexer := token.NewLexer(file)
}
