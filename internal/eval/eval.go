package eval

import (
	"fmt"

	"github.com/yurikdotdev/covfefescript/internal/ast"
	"github.com/yurikdotdev/covfefescript/internal/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		if node.Expression == nil {
			return object.COVFEFE
		}
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.LookStatement:
		return evalLookStatement(node, env)

	case *ast.MoneyLiteral:
		return &object.Money{Value: node.Value}
	case *ast.StringLiteral:
		return &object.Tweet{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.ForLoopStatement:
		return evalForLoopStatement(node, env)
	case *ast.BreakStatement:
		return object.BREAK
	case *ast.ContinueStatement:
		return object.CONTINUE
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	}

	return nil
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
