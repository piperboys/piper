package parser

import (
	"strconv"
)

type Expression interface {
	getType() string
}

type Integer struct {
	Value int
}

func (int Integer) getType() string {
	return "int"
}

type Operator struct {
	Value string
}

type Variable struct {
	Name  string
	Value Expression
}

func (variable Variable) getType() string {
	return variable.Value.getType()
}

type Operation struct {
	left     any
	operator Operator
	right    any
}

type LeftParenthesis struct {
	Value string
}

type RightParenthesis struct {
	Value string
}

type VariableDeclaration struct {
	Variable   Variable
	Expression any
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

func extractVariable(variableName string) (Variable, error) {
	return Variable{Name: variableName}, nil
}

func extractOperation(left any, operator any, right any) ([]any, error) {
	var combined []any

	leftSlice, isSlice := left.([]any)

	if isSlice {
		combined = append(combined, leftSlice...)
	} else {
		combined = append(combined, left)
	}

	combined = append(combined, operator.(Operator))

	rightSlice, isSlice := right.([]any)
	if isSlice {
		combined = append(combined, rightSlice...)
	} else {
		combined = append(combined, right)
	}

	return combined, nil
}

func extractGroup(expression any) ([]any, error) {
	combined := []any{LeftParenthesis{Value: "("}}

	sliceExpr, isSlice := expression.([]any)

	if isSlice {
		combined = append(combined, sliceExpr...)
	} else {
		combined = append(combined, expression)
	}

	combined = append(combined, RightParenthesis{Value: ")"})

	return combined, nil
}

func extractVariableDeclaration(variable any, expression any) (VariableDeclaration, error) {
	return VariableDeclaration{Variable: variable.(Variable), Expression: expression}, nil
}
