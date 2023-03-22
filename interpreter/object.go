package interpreter

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	IDENTIFIER_OBJ = "IDENTIFIER" // add, x, y, ...
	INTEGER_OBJ    = "INTEGER"    // 1343456
	STRING_OBJ     = "STRING"     // "foobar"
	BOOLEAN_OBJ    = "BOOLEAN"
	NUMBER_OBJ     = "NUMBER"
	FUNCTION_OBJ   = "FUNCTION"
	ARRAY_OBJ      = "ARRAY"
	BUILTIN_OBJ    = "BUILTIN"

	// compiler
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION_OBJ"
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

// ================= bool
type Boolean struct {
	Value bool
}

func (b Boolean) Type() string {
	return BOOLEAN_OBJ
}
func (b Boolean) String() string {
	res := "false"
	if b.Value {
		res = "true"
	}
	return res
}

// =================== builtin functions
type BuiltinFunction func(...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() string   { return BUILTIN_OBJ }
func (b *Builtin) String() string { return "builtin function" }

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

// ================== compiled function
type CompiledFunction struct {
	Instructions Instructions
	NumLocals    int
	NumParas     int
}

func (cf *CompiledFunction) Type() string {
	return COMPILED_FUNCTION_OBJ
}

func (cf *CompiledFunction) String() string {
	return fmt.Sprintf("func(%d paras, %d locals )", cf.NumParas, cf.NumLocals)
}

func randomObject(s string) Object {

	var res Object
	switch s {
	case "int":
		v := rand.Int() - 1
		if rand.Intn(2) == 0 { // [0, 2) = 0, 1
			v = -v
		}
		res = &Integer{Value: v}

	case "bool":
		res = &Boolean{Value: bool(rand.Intn(2) == 0)}

	case "string":
		res = &String{Value: " "}
	default:

		panic(fmt.Sprintf("invalid input = %s", s))
	}
	return res
}
