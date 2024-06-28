package interpreter

import "github.com/aerochrome/piper/internal/parser"

func TransformToReversePolishNotation(result []any) []any {
	if len(result) <= 1 {
		// TODO if function or possibly block, still order their expressions recursively and then return

		if len(result) == 1 {
			switch v := result[0].(type) {
			case parser.VariableDeclaration:
				// TODO Make this one better with pointers

				exprSlice, isSlice := v.Expression.([]any)
				if isSlice {
					v.Expression = TransformToReversePolishNotation(exprSlice)
				} else {
					v.Expression = TransformToReversePolishNotation([]any{v.Expression})
				}

				result[0] = v
			case parser.Function:
				exprSlice, isSlice := v.Expression.([]any)
				if isSlice {
					v.Expression = TransformToReversePolishNotation(exprSlice)
				} else {
					v.Expression = TransformToReversePolishNotation([]any{v.Expression})
				}

				result[0] = v
			}
		}

		// One element or less, it's already in the perfect order
		return result
	}

	var outputQueue []any
	var operatorStack []any

	for _, item := range result {
		switch v := item.(type) {
		case parser.Integer, parser.Variable:
			outputQueue = append(outputQueue, v)
		case parser.Operator:
			if len(operatorStack) > 0 {
				for idx := len(operatorStack) - 1; idx >= 0; idx-- { // Loop backwards (FIFO)
					operator, ok := operatorStack[idx].(parser.Operator)

					if !ok {
						_, isLeftParenthesis := operatorStack[idx].(parser.LeftParenthesis)

						if isLeftParenthesis {
							break
						}

						panic("Encountered unknown operator on operator stack")
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
		case parser.LeftParenthesis:
			operatorStack = append(operatorStack, v)
		case parser.RightParenthesis:
			matchedParenthesis := false

			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				_, isLeftParenthesis := top.(parser.LeftParenthesis)

				if isLeftParenthesis {
					operatorStack = operatorStack[:len(operatorStack)-1]
					matchedParenthesis = true
					break
				} else {
					outputQueue = append(outputQueue, top)
					operatorStack = operatorStack[:len(operatorStack)-1]
				}
			}

			if !matchedParenthesis {
				panic("Mismatched parenthesis: no matching opening parenthesis")
			}
		case parser.VariableDeclaration:
			// TODO Make this one better with pointers

			exprSlice, isSlice := v.Expression.([]any)
			if isSlice {
				v.Expression = TransformToReversePolishNotation(exprSlice)
			} else {
				v.Expression = TransformToReversePolishNotation([]any{v.Expression})
			}

			outputQueue = append(outputQueue, v)
		case parser.Function:
			exprSlice, isSlice := v.Expression.([]any)
			if isSlice {
				v.Expression = TransformToReversePolishNotation(exprSlice)
			} else {
				v.Expression = TransformToReversePolishNotation([]any{v.Expression})
			}

			outputQueue = append(outputQueue, v)
		default:
			panic("Unknown token encountered")
		}

	}

	// Move every operator left to outputQueue
	if len(operatorStack) > 0 {
		for idx := len(operatorStack) - 1; idx >= 0; idx-- {
			_, isLeftParenthesis := operatorStack[idx].(parser.LeftParenthesis)
			if isLeftParenthesis {
				panic("Mismatched parenthesis: no matching closing parenthesis")
			}

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
		"|":  1,
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
