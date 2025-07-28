package eval

import (
	"github.com/yurikdotdev/covfefescript/internal/ast"
	"github.com/yurikdotdev/covfefescript/internal/object"
)

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		if statement == nil {
			continue
		}
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.GiveMeValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	if result == nil {
		return object.COVFEFE
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()

			if rt == object.GIVE_ME_VALUE || rt == object.ERROR_OBJ || rt == "BREAK_VALUE" || rt == "CONTINUE_VALUE" {
				return result
			}
		}
	}

	return result
}

func evalReturnStatement(rs *ast.ReturnStatement, env *object.Environment) object.Object {
	val := Eval(rs.ReturnValue, env)
	if isError(val) {
		return val
	}
	return &object.GiveMeValue{Value: val}
}


func evalLookStatement(ls *ast.LookStatement, env *object.Environment) object.Object {
	val := Eval(ls.Value, env)
	if isError(val) {
		return val
	}
	env.Set(ls.Name.Value, val)
	return object.COVFEFE
}

func evalForLoopStatement(fl *ast.ForLoopStatement, env *object.Environment) object.Object {
	for {
		condition := Eval(fl.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break 
		}

		bodyResult := Eval(fl.Body, env)
		if bodyResult != nil {
			if bodyResult.Type() == "BREAK_VALUE" {
				break 
			}
			if bodyResult.Type() == "CONTINUE_VALUE" {
				continue 
			}
			if isError(bodyResult) {
				return bodyResult 
			}
		}
	}

	return object.COVFEFE
}