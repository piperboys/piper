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

func (interpreter *Interpreter) Evaluate(input []any) any {
	var stack []any

	for _, item := range input {
		switch item := item.(type) {
		case parser.Integer, parser.Function:
			stack = append(stack, item)
		case parser.Operator:
			left := stack[len(stack)-2]
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-2] // remove the last two

			stack = append(stack, evaluateOperation(left, right, item))
		case parser.VariableDeclaration:
			_, exists := interpreter.variables[item.Variable.Name]

			if exists {
				panic(fmt.Sprintf("Variable '%s' cannot be redefined", item.Variable.Name))
			}

			exprSlice, ok := item.Expression.([]any)

			var result any
			if ok {
				result = interpreter.Evaluate(exprSlice)
			} else {
				result = interpreter.Evaluate([]any{item.Expression})
			}

			interpreter.variables[item.Variable.Name] = &item.Variable
			interpreter.variables[item.Variable.Name].Value = result.(parser.Expression)

			stack = append(stack, result)
		case parser.Variable:
			variable, exists := interpreter.variables[item.Name]

			if !exists {
				panic(fmt.Sprintf("Variable '%s' is not defined", item.Name))
			}

			// TODO we might wanna push the variable on the stack instead of just the value?
			stack = append(stack, variable.Value)
		default:
			panic("Unknown token found")
		}
	}

	if len(stack) > 1 {
		panic("[Eval] Stack has more than one element")
	}

	return stack[0]
}

func evaluateOperation(left any, right any, operator parser.Operator) any {
	leftInt, ok1 := left.(parser.Integer)
	rightInt, ok2 := right.(parser.Integer)

	if !(ok1 && ok2) {
		panic("Invalid operation")
	}

	switch operator.Value {
	case "+":
		return parser.Integer{Value: leftInt.Value + rightInt.Value}
	case "-":
		return parser.Integer{Value: leftInt.Value - rightInt.Value}
	case "*":
		return parser.Integer{Value: leftInt.Value * rightInt.Value}
	case "/":
		return parser.Integer{Value: leftInt.Value / rightInt.Value}
	default:
		panic("Invalid operator")
	}
}
