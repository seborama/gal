package gal

type Variable struct {
	Name string
}

func NewVariable(name string) Variable {
	return Variable{
		Name: name,
	}
}

func (Variable) kind() entryKind {
	return variableEntryKind
}

func (v Variable) String() string {
	return string(v.Name)
}
