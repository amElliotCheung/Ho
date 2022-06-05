package interpreter

import (
	"log"
	"strings"
	"testing"
)

func TestLexer_Read(t *testing.T) {
	input := `a = 0
	b = "123"
	c = "asdf"
	d := "a\n\n\n"
	afunc := func(){}
	array := {1,b,c,d,afunc}
	array[2]
	1,2 -> 3
	hope`
	in := strings.NewReader(input)
	lexer := NewLexer(in)

	for tk := lexer.Read(); tk.Literal() != "EOF"; tk = lexer.Read() {
		log.Println(tk)
	}
}
