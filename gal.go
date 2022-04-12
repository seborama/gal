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
	variableType
	functionType
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

func parseParts(expr string) {
	for idx := 0; idx < len(expr); {
		part, pType, pos := extractPart(expr[idx:])
		if part == "" {
			break
		}

		if pos == 0 {
			fmt.Println(part)
			break
		}

		fmt.Printf("idx: %d >> type: %d >> part: '%s'\n", idx, pType, part)
		idx += pos
	}
}

func extractPart(expr string) (string, exprType, int) {
	// left trim blanks
	from := 0
	for _, r := range expr {
		if !isPartBoundary(r) {
			break
		}
		from++
	}

	// blank: no part
	if from == len(expr) {
		return "", blankType, from
	}

	// read part - "string"
	if expr[from] == '"' {
		s := readString(expr[from:])
		return s, stringType, len(s)
	}

	// read part - :variable:
	if expr[from] == ':' {
		s, l := readVariable(expr[from:])
		return s, variableType, l
	}

	// read part - function(...)
	// TODO: this is not complete: the (...) part of the function needs to be parsed and eval'ed
	if (expr[from] >= 'a' && expr[from] <= 'z') || (expr[from] >= 'A' && expr[from] <= 'Z') {
		s, l := readFunctionName(expr[from:])
		return s, functionType, l
	}

	// read part - operator
	// TODO: only single character operators are supported
	if isOperator(string(expr[from])) {
		if expr[from] == '+' || expr[from] == '-' {
			s, l := squashPlusMinusChain(expr[from:])
			return s, operatorType, l
		}
		return string(expr[from]), operatorType, 1
	}

	// read part - number
	// TODO: complex numbers are not supported
	to := 0
	isFloat := false
	for _, r := range expr[from:] {
		to++

		if isPartBoundary(r) || isOperator(string(r)) {
			break
		}

		if r == '.' && !isFloat {
			isFloat = true
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}

		return fmt.Sprintf("Syntax error: invalid character '%c' for number '%s'", r, expr[:to]), numericalType, 0
	}

	return expr[from:to], unknownType, to
}

func readString(expr string) string {
	to := 1 // keep leading double-quotes
	escapes := 0

	for i, r := range expr[1:] {
		to += 1
		if expr[i] == '\\' {
			escapes += 1
			continue
		}
		if r == '"' && (escapes == 0 || escapes&1 == 0) {
			break
		}

		escapes = 0
	}

	return expr[:to]
}

func readVariable(expr string) (string, int) {
	to := 1 // keep leading ':'

	for _, r := range expr[1:] {
		to += 1
		if r == ':' {
			break
		}
		if isPartBoundary(r) {
			return fmt.Sprintf("Syntax error: invalid character '%c' for variable name '%s'", r, expr[:to]), 0
		}
	}

	if expr[to-1] != ':' {
		return fmt.Sprintf("Syntax error: missing ':' to end variable '%s'", expr[:to]), 0
	}

	return expr[:to], len(expr[:to])
}

// f() is a function f with no args.
// f(x) is a function f with one arg 'x'.
// f(x, y, ...) is a function f with multiple args 'x', 'y', ...
// (...) is a standard associative grouping with parentheses. It is akin to an 'identity' function: `f(x) = x`.
// () is probably invalid
func readFunctionName(expr string) (string, int) {
	to := 1 // we know the first character is a letter, and it is the first letter

	for _, r := range expr[1:] {
		if r == '(' {
			return expr[:to], len(expr[:to])
		}
		if isPartBoundary(r) {
			return fmt.Sprintf("Syntax error: invalid character '%c' for function name '%s'", r, expr[:to]), 0
		}
		to += 1
	}

	return fmt.Sprintf("Syntax error: missing '(' for function name '%s'", expr[:to]), 0
}

func squashPlusMinusChain(expr string) (string, int) {
	to := 0
	outcomeSign := 1

	for _, r := range expr {
		if isPartBoundary(r) {
			break
		}
		if r != '+' && r != '-' {
			break
		}
		if r == '-' {
			outcomeSign = -outcomeSign
		}
		to += 1
	}

	sign := "-"
	if outcomeSign == 1 {
		sign = "+"
	}

	return sign, to
}

func isPartBoundary(r rune) bool {
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
