package gal

type Operator string

func (Operator) kind() entryKind {
	return operatorEntryKind
}

func (o Operator) String() string {
	return string(o)
}

const (
	invalidOperator    Operator = "invalid"
	Plus               Operator = "+"
	Minus              Operator = "-"
	Multiply           Operator = "*"
	Divide             Operator = "/"
	Modulus            Operator = "%"
	Power              Operator = "**"
	LShift             Operator = "<<"
	RShift             Operator = ">>"
	LessThan           Operator = "<"
	LessThanOrEqual    Operator = "<="
	EqualTo            Operator = "=="
	NotEqualTo         Operator = "!="
	GreaterThan        Operator = ">"
	GreaterThanOrEqual Operator = ">="
	And                Operator = "And" // NOTE: case sentive for now
	And2               Operator = "&&"
	Or                 Operator = "Or" // NOTE: case sentive for now
	Or2                Operator = "||"
)

func powerOperators(o Operator) bool {
	return o == Power
}

func multiplicativeOperators(o Operator) bool {
	return o == Multiply || o == Divide || o == Modulus
}

func additiveOperators(o Operator) bool {
	return o == Plus || o == Minus
}

func bitwiseShiftOperators(o Operator) bool {
	return o == LShift || o == RShift
}

func comparativeOperators(o Operator) bool {
	return o == GreaterThan || o == GreaterThanOrEqual ||
		o == LessThan || o == LessThanOrEqual ||
		o == EqualTo || o == NotEqualTo
}

func logicalOperators(o Operator) bool {
	return o == And || o == And2 ||
		o == Or || o == Or2
}
