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

func operatorPrecedence(o Operator) int {
	switch o {
	case invalidOperator:
		return 0
	case plus, minus:
		return 1
	case multiply, divide, modulus:
		return 2
	case power:
		return 3
	default:
		panic("unknown operator: " + o.String())
	}
}
