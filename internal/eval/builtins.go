package eval

import (
	"fmt"

	"github.com/yurikdotdev/covfefescript/internal/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Tweet:
				return &object.Money{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Money{Value: int64(len(arg.Elements))}
			default:
				return object.NewError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"BING": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return object.COVFEFE
		},
	},
	"SADLY": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.Error{Message: args[0].Inspect()}
		},
	},
}
