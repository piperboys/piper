package parser

import (
	"strconv"
)

type Integer struct {
	value int
}

type Operator struct {
	value string
}

type Operation struct {
	left     any
	operator Operator
	right    any
}

func extract_expression(input any) (any, error) {
	results := []any{}

	for _, item := range input.([]any) {
		lineResult := []any{}

		for _, token := range item.([]any)[0].([]any) {
			expression := token.([]any)[1]
			lineResult = append(lineResult, expression)
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

	return Operator{value: string(v)}, nil
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

	return Integer{value: value}, nil
}

func extractOperation(left any, operator any, right any) (Operation, error) {
	return Operation{
		left:     left,
		operator: operator.(Operator),
		right:    right,
	}, nil
}
