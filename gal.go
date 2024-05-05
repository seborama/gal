package gal

import "fmt"

type exprType int

const (
	unknownType exprType = iota
	blankType
	numericalType
	operatorType
	stringType
	variableType
	functionType
	boolType
)

type Value interface {
	// Calculation
	Add(Value) Value
	Sub(Value) Value
	Multiply(Value) Value
	Divide(Value) Value
	PowerOf(Value) Value
	Mod(Value) Value
	LShift(Value) Value
	RShift(Value) Value
	// Logical
	LessThan(Value) Bool
	LessThanOrEqual(Value) Bool
	EqualTo(Value) Bool
	NotEqualTo(Value) Bool
	GreaterThan(Value) Bool
	GreaterThanOrEqual(Value) Bool
	And(Value) Bool
	Or(Value) Bool
	// Helpers
	Stringer
	fmt.Stringer
	entry
}

// Example: Parse("blah").Eval(WithVariables(...), WithFunctions(...))
// This allows to parse an expression and then use the resulting Tree for multiple
// evaluations with different variables provided.
func Parse(expr string) Tree {
	treeBuilder := NewTreeBuilder()

	tree, err := treeBuilder.FromExpr(expr)
	if err != nil {
		return Tree{
			NewUndefinedWithReasonf(err.Error()),
		}
	}

	return tree
}
