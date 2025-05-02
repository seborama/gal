package gal

type Variable struct {
	Name string
}

func NewVariable(name string) Variable {
	return Variable{
		Name: name,
	}
}

func (v Variable) Calculate(val entry, op Operator, cfg *treeConfig) entry {
	varName := v.Name

	rhsVal := cfg.Variable(varName)
	if u, ok := rhsVal.(Undefined); ok {
		return u
	}

	if val == nil {
		return rhsVal
	}

	val = calculate(val.(Value), op, rhsVal)

	return val
}

func (v Variable) String() string {
	return v.Name
}
