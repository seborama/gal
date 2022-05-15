package gal

type exprType int

const (
	unknownType exprType = iota
	numericalType
	operatorType
	stringType
	variableType
	functionType
	blankType // TODO: remove since it's a non-expression?
)

type Value interface {
	Add(Value) Value
	Sub(Value) Value
	Multiply(Value) Value
	Divide(Value) Value
	PowerOf(Value) Value
	Mod(Value) Value
	Stringer
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
