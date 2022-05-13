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
	stringer
	entry
}

func Eval(expr string, opts ...option) Value {
	treeBuilder := NewTreeBuilder(opts...)

	tree, err := treeBuilder.FromExpr(expr)
	if err != nil {
		return NewUndefinedWithReasonf(err.Error())
	}

	return tree.Eval()
}
