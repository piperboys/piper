package parser

import (
	"errors"
	"strconv"
)

type Expression interface {
	GetType() string
}

type Integer struct {
	Value int
}

func (int Integer) GetType() string {
	return "int"
}

type Float64 struct {
	Value float64
}

func (float Float64) GetType() string {
	return "float"
}

type Operator struct {
	Value string
}

type Variable struct {
	Name  string
	Value Expression
}

func (variable Variable) GetType() string {
	return variable.Value.GetType()
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

type Function struct {
	ArgumentName      string
	ArgumentType      string
	ReturnType        string
	Expression        any
	AdditionalContext map[string]*Variable
}

func (function Function) GetType() string {
	// The function declaration in itself is of type func, only the call of the func has the ReturnType as Type
	return "func"
}

type Array struct {
	Expressions [][]any
}

func (array Array) GetType() string {
	return "array"
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

func extractFloat64(float any) (Float64, error) {
	value, err := strconv.ParseFloat(float.(string), 64)

	if err != nil {
		panic("Invalid float parsed!")
	}

	return Float64{Value: value}, nil
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

func extractFunction(argument any, argType any, returnType any, expr any) (Function, error) {
	return Function{
		ArgumentName:      argument.(Variable).Name,
		ArgumentType:      argType.(string),
		ReturnType:        returnType.(string),
		Expression:        expr,
		AdditionalContext: nil,
	}, nil
}

func extractArraySingle(expression any) (Array, error) {
	switch v := expression.(type) {
	case Expression:
		return Array{Expressions: [][]any{{v}}}, nil
	case []any:
		return Array{Expressions: [][]any{v}}, nil
	}

	return Array{}, errors.New("invalid expression")
}

func extractArray(expressions any) (Array, error) {
	return Array{Expressions: expressions.([][]any)}, nil
}

func extractExpressionList(expression any, expressionList any) ([][]any, error) {
	// TODO don't use append here (does dynamic alloc)

	var results [][]any

	switch v := expression.(type) {
	case Expression:
		results = append(results, []any{v})
	case []any:
		results = append(results, v)
	}

	switch v := expressionList.(type) {
	case Expression:
		results = append(results, []any{v})
	case []any:
		results = append(results, v)
	case [][]any:
		results = append(results, v...)
	}

	return results, nil
}
