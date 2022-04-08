package main

import (
	"log"
	"os"
	"regexp"
)

var regexPat = regexp.MustCompile(`\s*((//.*)|([0-9]+)|("(\\"|\\\\|\\n|[^"])*")|[A-Z_a-z][A-Z_a-z0-9]*|==|<=|>=|&&|\|\||[[:punct:]])?`)

func main() {
	text, err := os.ReadFile("./sourcecode.txt")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("File contents: %s", text)

}
