package eval

import (
	"github.com/yurikdotdev/covfefescript/internal/ast"
	"github.com/yurikdotdev/covfefescript/internal/object"
)

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	switch node.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return object.NewError("unknown operator: %s%s", node.Operator, right.Type())
	}
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	if node.Operator == "AND" {
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		if !isTruthy(left) {
			return left
		}

		return Eval(node.Right, env)
	}

	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.MONEY_OBJ && right.Type() == object.MONEY_OBJ:
		return evalMoneyInfixExpression(node.Operator, left, right)
	case left.Type() == object.TWEET_OBJ && right.Type() == object.TWEET_OBJ:
		return evalTweetInfixExpression(node.Operator, left, right)
	case left.Type() == object.TRUTH_OBJ && right.Type() == object.TRUTH_OBJ:
		switch node.Operator {
		case "==":
			return nativeBoolToBooleanObject(left == right)
		case "!=":
			return nativeBoolToBooleanObject(left != right)
		default:
			return object.NewError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
		}
	case left.Type() == object.COVFEFE_OBJ || right.Type() == object.COVFEFE_OBJ:
		return object.NewError("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	case node.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case node.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return object.NewError("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return object.COVFEFE
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return object.NewError("identifier not found: %s", node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.YUGE:
		return object.LOSER
	case object.LOSER:
		return object.YUGE
	case object.COVFEFE:
		return object.YUGE
	default:
		return object.LOSER
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.MONEY_OBJ {
		return object.NewError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Money).Value
	return &object.Money{Value: -value}
}

func evalMoneyInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Money).Value
	rightVal := right.(*object.Money).Value
	switch operator {
	case "+":
		return &object.Money{Value: leftVal + rightVal}
	case "-":
		return &object.Money{Value: leftVal - rightVal}
	case "*":
		return &object.Money{Value: leftVal * rightVal}
	case "/":
		return &object.Money{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "%":
		return &object.Money{Value: leftVal % rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalTweetInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return object.NewError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.Tweet).Value
	rightVal := right.(*object.Tweet).Value
	return &object.Tweet{Value: leftVal + rightVal}
}

func nativeBoolToBooleanObject(input bool) *object.Truth {
	if input {
		return object.YUGE
	}
	return object.LOSER
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.COVFEFE, object.LOSER:
		return false
	default:
		return true
	}
}
