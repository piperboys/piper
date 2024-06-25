package interpreter

import "github.com/aerochrome/piper/internal/parser"

func TransformToReversePolishNotation(result []any) []any {
	if len(result) <= 1 {
		// TODO if function or possibly block, still order their expressions recursively and then return

		// One element or less, it's already in the perfect order
		return result
	}

	var outputQueue []any
	var operatorStack []any

	for _, item := range result {
		switch v := item.(type) {
		case parser.Integer:
			outputQueue = append(outputQueue, v)
		case parser.Operator:
			if len(operatorStack) > 0 {
				for idx := len(operatorStack) - 1; idx >= 0; idx-- { // Loop backwards (FIFO)
					operator, ok := operatorStack[idx].(parser.Operator)

					if !ok {
						panic("Encountered non-operator on operator stack")
					}

					operatorPrecedence := getPrecedence(operator)
					valuePrecedence := getPrecedence(v)

					if (operatorPrecedence > valuePrecedence) || (operatorPrecedence == valuePrecedence && getAssociativity(v) == "left") {
						outputQueue = append(outputQueue, operator)
						operatorStack = operatorStack[:len(operatorStack)-1] // pop operator
					}
				}
			}

			operatorStack = append(operatorStack, v)
		}

	}

	// Move every operator left to outputQueue
	if len(operatorStack) > 0 {
		for idx := len(operatorStack) - 1; idx >= 0; idx-- {
			outputQueue = append(outputQueue, operatorStack[idx])
		}
	}

	return outputQueue
}

func getPrecedence(operator parser.Operator) int {
	precedenceMap := map[string]int{
		"*":  3,
		"/":  3,
		"+":  2,
		"-":  2,
		"|>": 1,
		"|":  0,
	}

	precedence, ok := precedenceMap[operator.Value]

	if !ok {
		panic("Missing precedence for operator " + operator.Value)
	}

	return precedence
}

func getAssociativity(operator parser.Operator) string {
	precedenceMap := map[string]string{
		"+":  "left",
		"-":  "left",
		"*":  "left",
		"/":  "left",
		"|>": "left",
		"|":  "left",
	}

	precedence, ok := precedenceMap[operator.Value]

	if !ok {
		panic("Missing associativity for operator " + operator.Value)
	}

	return precedence
}
