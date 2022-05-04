package interpreter

import (
	"bytes"
	"strconv"
	"strings"
)

const (
	IDENTIFIER_OBJ = "IDENTIFIER" // add, x, y, ...
	INTEGER_OBJ    = "INTEGER"    // 1343456
	STRING_OBJ     = "STRING"     // "foobar"
	NUMBER_OBJ     = "NUMBER"
	FUNCTION_OBJ   = "FUNCTION"
	ARRAY_OBJ      = "ARRAY"
)

type Object interface {
	Type() string
	String() string
}

// ================ integer
type Integer struct {
	Value int
}

func (i *Integer) Type() string {
	return INTEGER_OBJ
}
func (i *Integer) String() string {
	return strconv.Itoa(i.Value)
}

// ================ string
type String struct {
	Value string
}

func (s String) Type() string {
	return STRING_OBJ
}
func (s String) String() string {
	return s.Value
}

//
// =================== function object
//

type Function struct {
	Parameters []*IdentifierLiteral
	Body       *BlockExpression
	Env        *Environment
}

func (f *Function) Type() string { return FUNCTION_OBJ }
func (f *Function) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

//	============== Array

type Array struct {
	Elements []Object
}

func (a *Array) Type() string {
	return ARRAY_OBJ
}

func (a *Array) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.String())
	}
	out.WriteString(LBRACKET)
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString(RBRACKET)
	return out.String()
}
