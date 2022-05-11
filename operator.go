package gal

type Operator string

func (Operator) kind() entryKind {
	return operatorEntryKind
}

func (o Operator) String() string {
	return string(o)
}

const (
	invalidOperator Operator = "invalid"
	plus            Operator = "+"
	minus           Operator = "-"
	multiply        Operator = "*"
	divide          Operator = "/"
	modulus         Operator = "%"
	power           Operator = "^"
)

func powerOperators(o Operator) bool {
	return o == power
}

func multiplicativeOperators(o Operator) bool {
	return o == multiply || o == divide || o == modulus
}

func additiveOperators(o Operator) bool {
	return o == plus || o == minus
}
