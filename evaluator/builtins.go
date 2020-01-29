package evaluator

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/lxdlam/monkey-plus/lexer"
	"github.com/lxdlam/monkey-plus/object"
	"github.com/lxdlam/monkey-plus/parser"
	"io"
	"os"
	"strings"
)

var builtins map[string]*object.Builtin

func InitBuiltins() {
	builtins = map[string]*object.Builtin{
		"len": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				case *object.Hash:
					return &object.Integer{Value: int64(arg.Len())}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
		"first": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*object.Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}

				return NULL
			},
		},
		"last": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length > 0 {
					return arr.Elements[length-1]
				}

				return NULL
			},
		},
		"rest": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1, length-1)
					copy(newElements, arr.Elements[1:length])
					return &object.Array{Elements: newElements}
				}

				return NULL
			},
		},
		"push": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)

				newElements := make([]object.Object, length+1, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]

				return &object.Array{Elements: newElements}
			},
		},
		"set": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 3 {
					return newError("wrong number of arguments. got=%d, want=3", len(args))
				}

				if args[0].Type() != object.HASH_OBJ {
					return newError("argument to `set` must be HASH, got %s", args[0].Type())
				}

				if _, ok := args[1].(object.Hashable); !ok {
					return newError("unusable as hash key: %s", args[1].Type())
				}

				hash := args[0].(*object.Hash).Clone()
				hash.Set(args[1], args[2])

				return hash
			},
		},
		"contains": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				if args[0].Type() != object.HASH_OBJ {
					return newError("argument to `contains` must be HASH, got %s", args[0].Type())
				}

				if _, ok := args[1].(object.Hashable); !ok {
					return newError("unusable as hash key: %s", args[1].Type())
				}

				_, ok := args[0].(*object.Hash).Get(args[1])
				return nativeBoolToBooleanObject(ok)
			},
		},
		"delete": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				if args[0].Type() != object.HASH_OBJ {
					return newError("argument to `delete` must be HASH, got %s", args[0].Type())
				}

				if _, ok := args[1].(object.Hashable); !ok {
					return newError("unusable as hash key: %s", args[1].Type())
				}

				hash := args[0].(*object.Hash).Clone()
				hash.Delete(args[1])

				return hash
			},
		},
		"puts": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return NULL
			},
		},
		"eval": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != object.STRING_OBJ {
					return newError("argument to `eval` must be STRING, got %s", args[0].Type())
				}

				return runCodeInner(strings.NewReader(string(args[0].(*object.String).Value)), object.NewEnvironment())
			},
		},
		"load": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != object.STRING_OBJ {
					return newError("argument to `load` must be STRING, got %s", args[0].Type())
				}

				newEnv := object.NewEnvironment()
				path := string(args[0].(*object.String).Value)

				file, err := os.Open(path)
				if err != nil {
					return newError("load %s failed", path)
				}

				defer file.Close()

				result := runCodeInner(bufio.NewReader(file), newEnv)

				if isError(result) {
					return newError("load %s failed. Inner error is: %s", path, result.(*object.Error).Message)
				}

				env.Merge(newEnv)
				return TRUE
			},
		},
		"type": &object.Builtin{
			Fn: func(env *object.Environment, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				return object.NewStringObject(string(args[0].Type()))
			},
		},
	}
}

func runCodeInner(in io.Reader, env *object.Environment) object.Object {
	scanner := bufio.NewScanner(in)

	var codes bytes.Buffer

	for {
		scanned := scanner.Scan()

		if !scanned {
			break
		}
		codes.WriteString(scanner.Text() + "\n")
	}

	l := lexer.New(codes.String())
	p := parser.New(l)

	program := p.ParseProgram()

	errorLength := len(p.Errors())
	if errorLength != 0 {
		var errors []string
		for idx, msg := range p.Errors() {
			errors = append(errors, fmt.Sprintf(" parse error[%d/%d]: %s", idx, errorLength, msg))
		}
		return newError("parser error: %s", strings.Join(errors, "\n"))
	}

	evaluated := Eval(program, env)

	if evaluated != nil {
		return evaluated
	} else {
		return NULL
	}
}
