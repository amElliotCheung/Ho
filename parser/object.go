package interpreter

import "fmt"

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
