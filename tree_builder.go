package gal

import (
	"strings"

	"github.com/pkg/errors"
)

type TreeBuilder struct{}

func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{}
}

func (tb TreeBuilder) FromExpr(expr string) (Tree, error) {
	exprTree := Tree{}

	for idx := 0; idx < len(expr); {
		part, ptype, length, err := extractPart(expr[idx:])
		if err != nil {
			return nil, err
		}

		if ptype == blankType {
			break
		}

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
			case Plus.String():
				exprTree = append(exprTree, Plus)
			case Minus.String():
				exprTree = append(exprTree, Minus)
			case Multiply.String():
				exprTree = append(exprTree, Multiply)
			case Divide.String():
				exprTree = append(exprTree, Divide)
			case Modulus.String():
				exprTree = append(exprTree, Modulus)
			case Power.String():
				exprTree = append(exprTree, Power)
			default:
				return nil, errors.Errorf("unknown operator: '%s'", part)
			}

		case functionType:
			fname, l, _ := readFunctionName(part)
			v, err := tb.FromExpr(part[l+1 : len(part)-1]) // exclude leading '(' and trailing ')'
			if err != nil {
				return nil, err
			}
			if fname == "" {
				exprTree = append(exprTree, v)
			} else {
				bodyFn := PreDefinedFunction(fname)
				exprTree = append(exprTree, NewFunction(fname, bodyFn, v.Split()...))
			}

		case variableType:
			v := NewVariable(part)
			exprTree = append(exprTree, v)

		default:
			return nil, errors.Errorf("internal error: unknown expression part type '%v'", ptype)
		}

		idx += length
	}

	return exprTree, nil
}

// returns the part extracted as string, the type extracted, the cursor position
// after extraction or an error.
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
	if s, l := readOperator(expr[pos:]); l != 0 {
		if s == "+" || s == "-" {
			s, l = squashPlusMinusChain(expr[pos:]) // TODO: move this into readOperator()?
		}
		return s, operatorType, pos + l, nil
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
		return "", 0, errors.Errorf("syntax error: non-terminated string '%s'", expr[:to])
	}

	return expr[1 : to-1], to, nil
}

func readVariable(expr string) (string, int, error) {
	to := 1 // keep leading ':'

	for _, r := range expr[1:] {
		to += 1
		if r == ':' {
			break
		}
		if isBlankSpace(r) {
			return "", 0, errors.Errorf("syntax error: invalid character '%c' for variable name '%s'", r, expr[:to])
		}
	}

	if expr[to-1] != ':' {
		return "", 0, errors.Errorf("syntax error: missing ':' to end variable '%s'", expr[:to])
	}

	return expr[:to], to, nil
}

func readFunctionName(expr string) (string, int, error) {
	to := 0 // this could be an anonymous identity function (i.e. simple case of parenthesis grouping)

	for _, r := range expr {
		if r == '(' {
			return expr[:to], to, nil
		}
		if isBlankSpace(r) {
			return "", 0, errors.Errorf("syntax error: invalid character '%c' for function name '%s'", r, expr[:to])
		}
		to += 1
	}

	return "", 0, errors.Errorf("syntax error: missing '(' for function name '%s'", expr[:to])
}

func readFunctionArguments(expr string) (string, int, error) {
	to := 1
	bktCount := 1 // the currently opened bracket

	for i := 1; i < len(expr); i++ {
		r := expr[i]

		if r == '"' {
			_, l, err := readString(expr[to:])
			if err != nil {
				return "", 0, err
			}
			to += l
			i += l - 1
			continue
		}

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
	}

	return "", 0, errors.Errorf("syntax error: missing ')' for function arguments '%s'", expr[:to])
}

func readNumber(expr string) (string, int, error) {
	to := 0
	isFloat := false

	for i, r := range expr {
		if isBlankSpace(r) || isOperator(expr[i:]) {
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

		return "", 0, errors.Errorf("syntax error: invalid character '%c' for number '%s'", r, expr[:to])
	}

	return expr[:to], to, nil
}

func squashPlusMinusChain(expr string) (string, int) {
	to := 0
	outcomeSign := 1

	for _, r := range expr {
		// if isBlankSpace(r) {
		// 	break
		// }
		if r != '+' && r != '-' && !isBlankSpace(r) {
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

func readOperator(s string) (string, int) {
	if strings.HasPrefix(s, Power.String()) {
		return s[:2], 2
	}

	if strings.HasPrefix(s, Plus.String()) ||
		strings.HasPrefix(s, Minus.String()) ||
		strings.HasPrefix(s, Divide.String()) ||
		strings.HasPrefix(s, Multiply.String()) ||
		strings.HasPrefix(s, Power.String()) ||
		strings.HasPrefix(s, Modulus.String()) {
		return s[:1], 1
	}

	return "", 0
}

func isOperator(s string) bool {
	_, l := readOperator(s)
	return l != 0
}
