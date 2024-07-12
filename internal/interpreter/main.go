package interpreter

import (
	"fmt"
	"github.com/aerochrome/piper/internal/parser"
)

type Interpreter struct {
	variables map[string]*parser.Variable
}

func NewInterpreter() *Interpreter {
	return &Interpreter{variables: make(map[string]*parser.Variable)}
}

// Evaluate | additionalContext is for scope variables, like in a function call for example (the argument)
// the additionalContext also shadows the global variables
func (interpreter *Interpreter) Evaluate(input []any, additionalContext map[string]*parser.Variable) any {
	var stack []any

	for _, item := range input {
		switch item := item.(type) {
		case parser.Integer, parser.Function, parser.Float64:
			stack = append(stack, item)
		case parser.Operator:
			left := stack[len(stack)-2]
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-2] // remove the last two

			stack = append(stack, interpreter.evaluateOperation(left, right, item, additionalContext))
		case parser.VariableDeclaration:
			_, existsInContext := additionalContext[item.Variable.Name]
			_, existsInGlobal := interpreter.variables[item.Variable.Name]

			if existsInContext || existsInGlobal {
				panic(fmt.Sprintf("Variable '%s' cannot be redefined", item.Variable.Name))
			}

			exprSlice, ok := item.Expression.([]any)

			var result any
			if ok {
				result = interpreter.Evaluate(exprSlice, additionalContext)
			} else {
				result = interpreter.Evaluate([]any{item.Expression}, additionalContext)
			}

			interpreter.variables[item.Variable.Name] = &item.Variable
			interpreter.variables[item.Variable.Name].Value = result.(parser.Expression)

			stack = append(stack, result)
		case parser.Variable:
			contextVariable, existsInContext := additionalContext[item.Name]

			if existsInContext {
				stack = append(stack, contextVariable.Value)
			} else {
				variable, exists := interpreter.variables[item.Name]

				if !exists {
					panic(fmt.Sprintf("Variable '%s' is not defined", item.Name))
				}

				// TODO we might wanna push the variable on the stack instead of just the value?
				stack = append(stack, variable.Value)
			}
		case parser.Array:
			for idx, expr := range item.Expressions {
				item.Expressions[idx] = []any{interpreter.Evaluate(expr, additionalContext)}

				if item.ElementType == "" {
					item.ElementType = item.Expressions[idx][0].(parser.Expression).GetType()
				} else {
					if item.Expressions[idx][0].(parser.Expression).GetType() != item.ElementType {
						panic("Array expressions must have the same type")
					}
				}
			}

			stack = append(stack, item)
		default:
			panic("Unknown token found")
		}
	}

	if len(stack) > 1 {
		panic("[Eval] Stack has more than one element")
	}

	return stack[0]
}

func (interpreter *Interpreter) evaluateOperation(left any, right any, operator parser.Operator, additionalContext map[string]*parser.Variable) any {
	switch operator.Value {
	case "+", "-", "*", "/":
		return evaluateArithmeticOperation(left, right, operator)
	case "|":
		return interpreter.evaluateFunctionCall(left, right, additionalContext)
	case "|>":
		return interpreter.evaluateLoop(left, right, additionalContext)
	default:
		panic(fmt.Sprintf("Unknown operator '%s'", operator.Value))
	}
}

func (interpreter *Interpreter) evaluateLoop(left any, right any, additionalContext map[string]*parser.Variable) any {
	leftExpr, isArray := left.(parser.Array)

	if !isArray {
		panic(fmt.Sprintf("Left expression not an array, got '%T' instead", left))
	}

	switch v := right.(type) {
	case parser.Function:
		leftExpr.ElementType = ""

		for idx, item := range leftExpr.Expressions {
			leftExpr.Expressions[idx] = []any{interpreter.callFunction(v, item[0].(parser.Expression), additionalContext)}

			if leftExpr.ElementType == "" {
				leftExpr.ElementType = leftExpr.Expressions[idx][0].(parser.Expression).GetType()
			}

			// TODO we don't need to check here for different element types, because it comes from the same func with the same return type
		}

		return leftExpr
	default:
		panic(fmt.Sprintf("Expression '%v' of type '%T' is not a function", v, v))
	}
}

func (interpreter *Interpreter) evaluateFunctionCall(left any, right any, additionalContext map[string]*parser.Variable) any {
	leftExpr, isExpr := left.(parser.Expression)

	if !isExpr {
		panic(fmt.Sprintf("Unknown left expression type: %T", left))
	}

	switch v := right.(type) {
	case parser.Function:
		return interpreter.callFunction(v, leftExpr, additionalContext)
	default:
		panic(fmt.Sprintf("Expression '%v' of type '%T' is not a function", v, v))
	}
}

func (interpreter *Interpreter) callFunction(function parser.Function, argument parser.Expression, additionalContext map[string]*parser.Variable) any {
	if argument.GetType() != function.ArgumentType {
		panic(fmt.Sprintf("Cannot pass type '%s' to argument of type '%s'", argument.GetType(), function.ArgumentType))
	}

	functionContext := function.AdditionalContext

	if functionContext == nil {
		functionContext = make(map[string]*parser.Variable)
	} else {
		functionContext = copyAdditionalContextShallow(functionContext)
	}

	if additionalContext != nil {
		for key, value := range additionalContext {
			functionContext[key] = value
		}
	}

	functionContext[function.ArgumentName] = &parser.Variable{Name: function.ArgumentName, Value: argument}

	exprSlice, ok := function.Expression.([]any)

	var result any
	if ok {
		result = interpreter.Evaluate(exprSlice, functionContext)
	} else {
		result = interpreter.Evaluate([]any{function.Expression}, functionContext)
	}

	if result.(parser.Expression).GetType() != function.ReturnType {
		panic(fmt.Sprintf("Function returns '%s' but return type is '%s'", result.(parser.Expression).GetType(), function.ReturnType))
	}

	if result.(parser.Expression).GetType() == "func" {
		resultFunc := result.(parser.Function)
		resultFunc.AdditionalContext = functionContext

		result = resultFunc
	}

	return result
}

func evaluateArithmeticOperation(left any, right any, operator parser.Operator) any {
	// TODO refactor this monstrosity (works for my low 5am standards, but defo should be simplified somehow)

	isFloatOperation := false
	var leftValue any
	var rightValue any

	switch v := left.(type) {
	case parser.Integer:
		leftValue = v.Value
	case parser.Float64:
		isFloatOperation = true
		leftValue = v.Value
	}

	switch v := right.(type) {
	case parser.Integer:
		if isFloatOperation {
			rightValue = float64(v.Value)
		} else {
			rightValue = v.Value
		}
	case parser.Float64:
		if !isFloatOperation {
			isFloatOperation = true
			leftValue = float64(leftValue.(int))
		}
		rightValue = v.Value
	}

	var resultValue any

	switch operator.Value {
	case "+":
		switch v := leftValue.(type) {
		case int:
			resultValue = v + rightValue.(int)
		case float64:
			resultValue = v + rightValue.(float64)
		}
	case "-":
		switch v := leftValue.(type) {
		case int:
			resultValue = v - rightValue.(int)
		case float64:
			resultValue = v - rightValue.(float64)
		}
	case "*":
		switch v := leftValue.(type) {
		case int:
			resultValue = v * rightValue.(int)
		case float64:
			resultValue = v * rightValue.(float64)
		}
	case "/":
		switch leftValue.(type) {
		case int:
			resultValue = float64(leftValue.(int)) / float64(rightValue.(int))
		case float64:
			resultValue = leftValue.(float64) / rightValue.(float64)
		}
	default:
		panic("Invalid arithmetic operator")
	}

	switch v := resultValue.(type) {
	case int:
		return parser.Integer{Value: v}
	case float64:
		return parser.Float64{Value: v}
	default:
		panic("Congratz, you've found the impossible bug")
	}
}

// Function to make a shallow copy of a map
func copyAdditionalContextShallow(original map[string]*parser.Variable) map[string]*parser.Variable {
	mapCopy := make(map[string]*parser.Variable)
	for key, value := range original {
		mapCopy[key] = value
	}
	return mapCopy
}
