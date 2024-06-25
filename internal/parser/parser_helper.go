package parser

import (
	"strconv"
)

type Integer struct {
	Value int
}

type Operator struct {
	Value string
}

type Operation struct {
	left     any
	operator Operator
	right    any
}

func extractExpression(input any) (any, error) {
	var results []any

	for _, line := range input.([]any) {
		var lineResult []any

		// Get the second element in the line array (the expression)
		lineSlice := line.([]any)

		expression, ok := lineSlice[1].([]any)

		if !ok {
			lineResult = append(lineResult, line.([]any)[1])
		} else {
			lineResult = append(lineResult, expression...)
		}

		results = append(results, lineResult)
	}

	return results, nil
}

func extractOperator(operator any) (Operator, error) {
	v, ok := operator.([]uint8)

	if !ok {
		panic("Invalid operator parsed")
	}

	return Operator{Value: string(v)}, nil
}

func extractInteger(integer any) (Integer, error) {
	valueStr := ""

	for _, item := range integer.([]any) {
		valueStr += string(item.([]uint8))
	}

	value, err := strconv.Atoi(valueStr)

	if err != nil {
		panic("Invalid integer parsed!")
	}

	return Integer{Value: value}, nil
}

func extractOperation(left any, operator any, right any) ([]any, error) {
	var combined []any

	combined = append(combined, left)
	combined = append(combined, operator.(Operator))

	rightSlice, ok := right.([]any)
	if !ok {
		combined = append(combined, right)
	} else {
		combined = append(combined, rightSlice...)
	}

	return combined, nil
}
