package interpreter

import "log"

var builtins = map[string]*Builtin{
	"len": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				log.Panic("incorrect number of arguments in len function")
			}
			arg := args[0]
			switch arg := arg.(type) {
			case *String:
				return &Integer{Value: len(arg.Value)}
			case *Array:
				return &Integer{Value: len(arg.Elements)}
			default:
				log.Panic("wrong argument type in len function")
			}
			return &Integer{Value: 0}
		},
	},
	"append": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				log.Panic("incorrect number of arguments in push function")
			}
			log.Printf("---------- append ----------")
			log.Printf("%T %v %T %v", args[0], args[0], args[1], args[1])
			arr := args[0].(*Array)
			arr.Elements = append(arr.Elements, args[1])
			return arr
		},
	},
}
