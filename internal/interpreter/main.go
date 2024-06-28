package interpreter

import "github.com/aerochrome/piper/internal/parser"

func Evaluate(input []any) any {
	var stack []any

	for _, item := range input {
		switch item := item.(type) {
		case parser.Integer:
			stack = append(stack, item)
		case parser.Operator:
			left := stack[len(stack)-2]
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-2] // remove the last two

			stack = append(stack, evaluateOperation(left, right, item))
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
