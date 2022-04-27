package gal

import (
	"fmt"

	"github.com/pkg/errors"
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

func (Operator) kind() entryKind {
	return operatorEntryKind
}

func (o Operator) String() string {
	return string(o)
}

type Function struct {
	Name string
	Tree Tree
}

func NewFunction(name string, tree Tree) *Function {
	return &Function{
		Name: name,
		Tree: tree,
	}
}

func (Function) kind() entryKind {
	return functionEntryKind
}

func (f Function) String() string {
	return string(f.Name)
}

type entryKind int

const (
	unknownEntryKind entryKind = iota
	valueEntryKind
	operatorEntryKind
	treeEntryKind
	functionEntryKind
)

type entry interface {
	kind() entryKind
}

const (
	invalidOperator Operator = ""
	plus            Operator = "+"
	minus           Operator = "-"
	times           Operator = "*"
	dividedBy       Operator = "/"
	modulus         Operator = "%"
)

func operatorPrecedence(o Operator) int {
	switch o {
	case plus, minus:
		return 0
	case times, dividedBy, modulus:
		return 1
	default:
		panic("unknown operator: " + o.String())
	}
}

func Eval(expr string) Value {
	panic("TODO")
}

// TODO: perhaps return []Value rather than Value
func parseParts(expr string) (Value, error) {
	tree, err := buildExprTree(expr)
	if err != nil {
		return nil, err
	}

	value := tree.PrioritiseOperators().Eval()

	return value, nil
}

func buildExprTree(expr string) (Tree, error) {
	exprTree := Tree{}

	for idx := 0; idx < len(expr); {
		part, ptype, length, err := extractPart(expr[idx:])
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		if ptype == blankType {
			break
		}

		fmt.Printf("idx: %d >> type: %d >> part: '%s'\n", idx, ptype, part)

		switch ptype {
		case numericalType:
			v, err := NewNumberFromString(part)
			if err != nil {
				return nil, err
			}
			exprTree = append(exprTree, v)

		case stringType:
			v := NewString(part)
			exprTree = append(exprTree, v)

		case operatorType:
			switch part {
			case plus.String():
				exprTree = append(exprTree, plus)
			case minus.String():
				exprTree = append(exprTree, minus)
			case times.String():
				exprTree = append(exprTree, times)
			case dividedBy.String():
				exprTree = append(exprTree, dividedBy)
			case modulus.String():
				exprTree = append(exprTree, modulus)
			default:
				return nil, errors.WithStack(newErrUnknownOperator(part))
			}

		case functionType:
			// TODO: squash the leading and trailing '()'
			fname, l, _ := readFunctionName(part)
			fmt.Printf("  >>> sub: fname: '%s'\n", fname)
			v, err := buildExprTree(part[l+1 : len(part)-1]) // exclude leading '(' and trailing ')'
			fmt.Printf("  <<< sub: fname: '%s' >> err: '%v'\n", fname, err)
			if err != nil {
				return nil, err
			}
			if fname == "" {
				exprTree = append(exprTree, v)
			} else {
				exprTree = append(exprTree, NewFunction(fname, v))
			}
		}

		idx += length
	}

	return exprTree, nil
}

func extractPart(expr string) (string, exprType, int, error) {
	// left trim blanks
	pos := 0
	for _, r := range expr {
		if !isBlankSpace(r) {
			break
		}
		pos++
	}

	// blank: no part
	if pos == len(expr) {
		return "", blankType, pos, nil
	}

	// read part - "string"
	if expr[pos] == '"' {
		s, l, err := readString(expr[pos:])
		if err != nil {
			return "", unknownType, 0, err
		}
		return s, stringType, pos + l, nil
	}

	// read part - :variable:
	if expr[pos] == ':' {
		s, l, err := readVariable(expr[pos:])
		if err != nil {
			return "", unknownType, 0, err
		}
		return s, variableType, pos + l, nil
	}

	// read part - function(...)
	// conceptually, parenthesis grouping is a special case of anonymous identity function
	if expr[pos] == '(' || (expr[pos] >= 'a' && expr[pos] <= 'z') || (expr[pos] >= 'A' && expr[pos] <= 'Z') {
		fname, lf, err := readFunctionName(expr[pos:])
		if err != nil {
			return "", unknownType, 0, err
		}
		fargs, la, err := readFunctionArguments(expr[pos+lf:])
		if err != nil {
			return "", unknownType, 0, err
		}
		return fname + fargs, functionType, pos + lf + la, nil
	}

	// read part - operator
	// TODO: only single character operators are supported
	if isOperator(string(expr[pos])) {
		if expr[pos] == '+' || expr[pos] == '-' {
			s, l := squashPlusMinusChain(expr[pos:])
			return s, operatorType, pos + l, nil
		}
		return string(expr[pos]), operatorType, pos + 1, nil
	}

	// read part - number
	// TODO: complex numbers are not supported
	s, l, err := readNumber(expr[pos:])
	if err != nil {
		return "", unknownType, 0, err
	}
	return s, numericalType, pos + l, nil
}

func readString(expr string) (string, int, error) {
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

	if expr[to-1] != '"' {
		return "", 0, errors.WithStack(newErrSyntaxError(fmt.Sprintf("non-terminated string '%s'", expr[:to])))
	}

	return expr[:to], to, nil
}

func readVariable(expr string) (string, int, error) {
	to := 1 // keep leading ':'

	for _, r := range expr[1:] {
		to += 1
		if r == ':' {
			break
		}
		if isBlankSpace(r) {
			return "", 0, errors.WithStack(newErrSyntaxError(fmt.Sprintf("invalid character '%c' for variable name '%s'", r, expr[:to])))
		}
	}

	if expr[to-1] != ':' {
		return "", 0, errors.WithStack(newErrSyntaxError(fmt.Sprintf("missing ':' to end variable '%s'", expr[:to])))
	}

	return expr[:to], to, nil
}

// f() is a function f with no args.
// f(x) is a function f with one arg 'x'.
// f(x, y, ...) is a function f with multiple args 'x', 'y', ...
// (...) is a standard associative grouping with parentheses. It is akin to an 'identity' function: `f(x) = x`.
// () just does nothing more than waste space
func readFunctionName(expr string) (string, int, error) {
	to := 0 // this could be an anonymous identity function (i.e. simple case of parenthesis grouping)

	for _, r := range expr {
		if r == '(' {
			return expr[:to], to, nil
		}
		if isBlankSpace(r) {
			return "", 0, errors.WithStack(newErrSyntaxError(fmt.Sprintf("invalid character '%c' for function name '%s'", r, expr[:to])))
		}
		to += 1
	}

	return "", 0, errors.WithStack(newErrSyntaxError(fmt.Sprintf("missing '(' for function name '%s'", expr[:to])))
}

func readFunctionArguments(expr string) (string, int, error) {
	to := 1
	bktCount := 1 // the currently opened bracket

	for _, r := range expr[1:] {
		to += 1
		if r == '(' {
			bktCount++
			continue
		}
		if r == ')' {
			bktCount--
			if bktCount == 0 {
				return expr[:to], to, nil
			}
		}
		// TODO: handle stringType
	}

	return "", 0, errors.WithStack(newErrSyntaxError(fmt.Sprintf("missing ')' for function arguments '%s'", expr[:to])))
}

func readNumber(expr string) (string, int, error) {
	to := 0
	isFloat := false

	for _, r := range expr {
		if isBlankSpace(r) || isOperator(string(r)) {
			break
		}

		to++

		if r == '.' && !isFloat {
			isFloat = true
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}

		return "", 0, errors.WithStack(newErrSyntaxError(fmt.Sprintf("invalid character '%c' for number '%s'", r, expr[:to])))
	}

	return expr[:to], to, nil
}

func squashPlusMinusChain(expr string) (string, int) {
	to := 0
	outcomeSign := 1

	for _, r := range expr {
		if isBlankSpace(r) {
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

func isBlankSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "/" || s == "*" || s == "^" || s == "%"
}
