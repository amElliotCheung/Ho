package interpreter

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

type Object interface {
	Type() string
	Inspect() string
	GetValue() int
}

// ================ integer
type Integer struct {
	Value int
}

func (i *Integer) Type() string {
	return "INTEGER"
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) GetValue() int {
	return i.Value
}

//
//				identifier
//
// type Identifier struct {
// 	Name string
// }

// func (i *Identifier) Type() string {
// 	return "IDENTIFIER"
// }
// func (i *Identifier) Inspect() string {
// 	return i.Name
// }

// func (i *Identifier) GetValue() int {
// 	return 0
// }

// ====================== variable environment

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}
func (e *Environment) Get(name string) Object {
	defer func() {
		log.Println("the environment is ------->\t")
		for k, v := range e.store {
			log.Println(k, v)
		}
		log.Println("")
	}()

	obj, ok := e.store[name]
	if !ok {
		defaultObject := &Integer{Value: 0}
		e.store[name] = defaultObject
		return defaultObject
	}
	return obj
}

func (e *Environment) Set(name string, val Object) {
	e.store[name] = val
}

//
// =================== function object
//

type Function struct {
	parameters []*IdToken
	body       ASTNode
}

func (f *Function) Type() string {
	return "FUNCTION"
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.parameters {
		params = append(params, p.GetText())
	}
	out.WriteString("func (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.body.String())
	out.WriteString("\n}")

	return out.String()
}
