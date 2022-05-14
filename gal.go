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

// Example: Parse("blah", WithFunctions(...)).Eval(WithVariables(...))
// This allows to parse an expression and then use the resulting Tree for multiple
// evaluations with different variables provided.
func Parse(expr string, opts ...parseOption) Tree {
	treeBuilder := NewTreeBuilder(opts...)

	tree, err := treeBuilder.FromExpr(expr)
	if err != nil {
		return Tree{
			NewUndefinedWithReasonf(err.Error()),
		}
	}

	return tree
}

type Functions map[string]FunctionalValue

func (f Functions) Function(name string) (FunctionalValue, bool) {
	if f == nil {
		return nil, false
	}

	val, ok := f[name]
	return val, ok
}

type parseConfig struct {
	functions Functions
}

type parseOption func(*parseConfig)

func WithFunctions(f Functions) parseOption {
	return func(cfg *parseConfig) {
		cfg.functions = f
	}
}
