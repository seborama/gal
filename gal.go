package gal

type exprType int

const (
	unknownType exprType = iota
	blankType
	numericalType
	operatorType
	stringType
	boolType
	variableType
	functionType
	objectPropertyType           // "cousin" of a variableType, but for a property of a user-defined object
	objectMethodType             // "cousin" of a functionType, but for a method of a user-defined object
	objectAccessorByPropertyType // represents an object accessor of a "left hand side" expression by property
	objectAccessorByMethodType   // represents an object accessor of a "left hand side" expression by method
)

// Example: Parse("blah").Eval(WithVariables(...), WithFunctions(...), WithObjects(...))
// This allows to parse an expression and then use the resulting Tree for multiple
// evaluations with different variables provided.
func Parse(expr string) Tree {
	treeBuilder := NewTreeBuilder()

	tree, err := treeBuilder.FromExpr(expr)
	if err != nil {
		return Tree{
			NewUndefinedWithReasonf("%s", err.Error()),
		}
	}

	return tree
}
