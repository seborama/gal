package gal

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type exprType int

const (
	unknownType exprType = iota
	numericalType
	operatorType
	stringType
	blankType // TODO: remove
)

type stringer interface {
	String() string
}

type numberer interface {
	Number() decimal.Decimal
}

type Value interface {
	Add(Value) Value
	stringer
}

type Operator string

const (
	plus Operator = "+"
)

type String struct {
	value string
}

func NewString(s string) String {
	return String{
		value: s,
	}
}

func (s String) Add(other Value) Value {
	switch v := other.(type) {
	case String:
		return String{
			value: s.value + v.value,
		}
	}

	v, ok := other.(stringer)
	if !ok {
		return Undefined{}
	}

	return String{
		value: s.value + v.String(),
	}
}

func (s String) String() string {
	return s.value
}

type Number struct {
	value decimal.Decimal
}

func NewNumberFromString(s string) (Number, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Number{}, err
	}

	return Number{
		value: d,
	}, nil
}

func (n Number) Add(other Value) Value {
	switch v := other.(type) {
	case Number:
		return Number{
			value: n.value.Add(v.value),
		}
	}

	v, ok := other.(numberer)
	if !ok {
		return Undefined{}
	}

	return Number{
		value: n.value.Add(v.Number()),
	}
}

func (n Number) String() string {
	return n.value.String()
}

type Undefined struct{}

func (u Undefined) Add(Value) Value {
	return Undefined{}
}

func (u Undefined) String() string {
	return "undefined"
}

func Eval(expr string) Value {
	v := eval(expr)
	fmt.Printf("result value: '%+s'\n", v.String())
	return v
}

func eval(expr string) Value {
	var v Value

	length := len(expr)

	for i := 0; i < length; i++ {
		j := 0
		if i == 0 { // Yuk: move that out of the loop and refactor
			v, j = value(expr[i:])
			i += j
		}

		var o Operator
		o, j = operator(expr[i:])
		i += j

		v, j = operate(o, v, expr[i:])
		i += j
	}

	return v
}

func value(expr string) (Value, int) {
	// part := extractPart(expr)

	// if part[0] == '"' {
	// 	return NewString(part), len(part)
	// }

	// if v, err := NewNumberFromString(part); err == nil {
	// 	return v, len(part)
	// }

	panic("should never reached this point")
}

func operator(expr string) (Operator, int) {
	length := len(expr)

	for i := 0; i < length; i++ {
		// part := extractPart(expr[i:])
		// i += len(part)

		// switch partType {
		// case operatorType:
		// 	switch part {
		// 	case "+":
		// 		return plus, i
		// 	}
		// 	panic("should never reached this point")
		// }
	}

	panic("should never reached this point")
}

func partType(expr string) exprType {
	if expr[0] == '"' && expr[len(expr)-1] == '"' {
		return stringType
	}

	if _, err := decimal.NewFromString(expr); err == nil {
		return numericalType
	}

	if isOperator(expr) {
		return operatorType
	}

	return unknownType
}

func extractPart(expr string) (string, int) {
	// left trim blanks
	from := 0
	for _, r := range expr {
		if !isValueBoundary(r) {
			break
		}
		from++
	}

	if from == len(expr) {
		return "", from
	}

	// read part
	// if expr[from] == '"' {
	// 	return readString(expr[from:])
	// }

	to := from
	newFrom := from
	if expr[from] == '+' || expr[from] == '-' {
		newFrom++
		to++
	}

	for _, r := range expr[newFrom:] {
		if isValueBoundary(r) || isOperator(string(r)) {
			break
		}
		to++
	}

	return expr[from:to], to
}

func readString(expr string) string {
	s := "\""
	to := 0
	escapes := 0

	for i, r := range expr[1:] {
		to += 1
		if expr[i-1] == '\\' {
			escapes += 1
			continue
		}
		if r == '"' && (escapes == 0 || escapes&1 == 0) {
			break
		}

		escapes = 0
	}

	return s[:to]
}

func isValueBoundary(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' ||
		r == '(' || r == ')'
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "/" || s == "*" || s == "^" || s == "%"
}

func getValueType(c rune) exprType {
	switch {
	case c == ' ', c == '\t', c == '\n':
		return blankType

	case c >= '0' && c <= '9',
		c == '-',
		c == '_',
		c == '.':
		return numericalType

	case c >= 'a' && c <= 'z',
		c >= 'A' && c <= 'Z',
		c == '_':
		return stringType

	default:
		return unknownType
	}
}

func getOperatorType(c rune) exprType {
	switch {
	case c == ' ', c == '\t', c == '\n':
		return blankType

	case c >= '0' && c <= '9',
		c == '_',
		c == '.':
		return numericalType

	case c >= 'a' && c <= 'z',
		c >= 'A' && c <= 'Z',
		c == '_':
		return stringType

	case c == '+',
		c == '-',
		c == '*',
		c == '/',
		c == '^',
		c == '%':
		return operatorType

	default:
		return unknownType
	}
}

func operate(op Operator, lhsValue Value, expr string) (Value, int) {
	rhsValue, i := value(expr)

	if op == plus {
		return add(lhsValue, rhsValue), i
	}

	panic("not implemented yet")
}

func add(a, b Value) Value {
	return a.Add(b)
}
