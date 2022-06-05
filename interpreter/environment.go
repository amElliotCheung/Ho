package interpreter

import "fmt"

// ============== environment
type Environment struct {
	store  map[string]Object
	outter *Environment
}

func NewEnvironment(outter *Environment) *Environment {
	return &Environment{
		store:  make(map[string]Object),
		outter: outter,
	}
}

func NewFunctionEnvirontment(outter *Environment, fn *Function, args []Object) *Environment {
	env := NewEnvironment(outter)
	for k, v := range fn.Env.store {
		env.Set(k, v)
	}

	for i, p := range fn.Parameters {
		env.Set(p.Key, args[i])
	}
	return env
}
func (e *Environment) Get(name string) (Object, error) {
	cur := e // must exists
	for cur != nil {
		obj, ok := cur.store[name]
		if ok {
			return obj, nil
		}

		cur = cur.outter

	}
	return nil, fmt.Errorf("identifier doesn't exist")
}

func (e *Environment) Set(name string, val Object) {
	e.store[name] = val
}
