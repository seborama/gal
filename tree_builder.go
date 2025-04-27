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
	tree := Tree{}

	//nolint:errcheck // life's too short to check for type assertion success here
	for idx := 0; idx < len(expr); {
		part, ptype, length, err := extractPart(expr[idx:])
		if err != nil {
			return nil, err
		}

		switch ptype {
		case numericalType:
			v, err := NewNumberFromString(part)
			if err != nil {
				return nil, err
			}
			tree = append(tree, v)

		case stringType:
			v := NewString(part)
			tree = append(tree, v)

		case boolType:
			v, err := NewBoolFromString(part)
			if err != nil {
				return nil, err
			}
			tree = append(tree, v)

		case operatorType:
			switch part {
			case Plus.String():
				tree = append(tree, Plus)
			case Minus.String():
				tree = append(tree, Minus)
			case Multiply.String():
				tree = append(tree, Multiply)
			case Divide.String():
				tree = append(tree, Divide)
			case Modulus.String():
				tree = append(tree, Modulus)
			case Power.String():
				tree = append(tree, Power)
			case LessThan.String():
				tree = append(tree, LessThan)
			case LessThanOrEqual.String():
				tree = append(tree, LessThanOrEqual)
			case EqualTo.String():
				tree = append(tree, EqualTo)
			case NotEqualTo.String():
				tree = append(tree, NotEqualTo)
			case GreaterThan.String():
				tree = append(tree, GreaterThan)
			case GreaterThanOrEqual.String():
				tree = append(tree, GreaterThanOrEqual)
			case LShift.String():
				tree = append(tree, LShift)
			case RShift.String():
				tree = append(tree, RShift)
			case And.String():
				tree = append(tree, And)
			case And2.String():
				tree = append(tree, And2) // NOTE: re-route to And?
			case Or.String():
				tree = append(tree, Or)
			case Or2.String():
				tree = append(tree, Or2) // NOTE: re-route to Or?
			default:
				return nil, errors.Errorf("unknown operator: '%s'", part)
			}

		case functionType:
			fname, l, _ := readNamedExpressionType(part)   //nolint:errcheck // ignore err: we already parsed the function name when in extractPart()
			v, err := tb.FromExpr(part[l+1 : len(part)-1]) // parse the function's arguments: exclude leading '(' and trailing ')'
			if err != nil {
				return nil, err
			}
			if fname == "" {
				// parenthesis grouping, not a real function per-se.
				// conceptually, parenthesis grouping is a special case of anonymous identity function
				tree = append(tree, v)
			} else {
				bodyFn := BuiltInFunction(fname) // will be nil if it isn't a built-in function (i.e. user-defined or object method)
				// TODO: if bodyFn == nil, we have either a user-defined function or an user-defined object method. It may be worth
				// ...   creating a new type for this case. This could simplify the code by keeping each case separate and simpler.
				tree = append(tree, NewFunction(fname, bodyFn, v.Split()...))
			}

		case objectMethodType:
			// an objectMethodType represents a method access on a user-defined object
			fname, l, _ := readNamedExpressionType(part)   //nolint:errcheck // ignore err: we already parsed the function name when in extractPart()
			v, err := tb.FromExpr(part[l+1 : len(part)-1]) // parse the function's arguments: exclude leading '(' and trailing ')'
			if err != nil {
				return nil, err
			}
			splits := strings.SplitN(fname, ".", 2) // there should only ever be exactly 2 parts at this point
			om := NewObjectMethod(splits[0], splits[1], v.Split()...)
			tree = append(tree, om)

		case variableType:
			v := NewVariable(part)
			tree = append(tree, v)

		case objectPropertyType:
			// an objectPropertyType represents a property access on a user-defined object
			splits := strings.SplitN(part, ".", 2) // there should only ever be exactly 2 parts at this point
			v := NewObjectProperty(splits[0], splits[1])
			tree = append(tree, v)

		case objectAccessorByPropertyType:
			// an objectAccessorByPropertyType is an access to a property of an object retrieved from the last expression evaluated in the Tree.
			tree = append(tree, Dot[Variable]{
				Member: NewVariable(part[1:]), // skip the "."
			})

		case objectAccessorByMethodType:
			// an objectAccessorByMethodType is an access to a method of an object retrieved from the last expression evaluated in the Tree.
			v, err := tb.FromExpr(part[1:]) // skip the "."
			if err != nil {
				return nil, err
			}
			if len(v) != 1 {
				// NOTE: this should never happen because we have already extracted the object accessor into a single "part".
				return nil, errors.Errorf("syntax error: invalid object accessor function: '%s'", part)
			}
			m := v[0].(Function)
			if m.BodyFn != nil {
				// NOTE: this could be supported but it would turn the object into a prototype model e.g. like JavaScript
				return nil, errors.Errorf("internal error: invalid object accessor function: '%s' - BodyFn is not empty: this indicates the object's method was confused for a build-in function", part)
			}
			tree = append(tree, Dot[Function]{Member: m})

		case blankType:
			// only returned when the entire expression is empty or only contains blanks.
			return tree, nil

		default:
			return nil, errors.Errorf("internal error: unknown expression part type '%T'='%v'", ptype, ptype)
		}

		idx += length
	}

	// adjust trees that start with "Plus" or "Minus" followed by a "Numberer"
	if tree.TrunkLen() >= 2 {
		switch tree[0] {
		case Plus:
			return tree[1:], nil
		case Minus:
			return append(Tree{NewNumberFromInt(-1), Multiply}, tree[1:]...), nil
		}
	}

	return tree, nil
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

	// read part - constants
	// e.g. Phi (golden ratio), etc, user-defined or built-in (True, False for booleans)
	if s, l, ctype, ok := readConstant(expr[pos:]); ok {
		return s, ctype, pos + l, nil
	}

	// read part - :variable:
	if expr[pos] == ':' {
		s, l, err := readVariable(expr[pos:])
		if err != nil {
			return "", unknownType, 0, err
		}
		return s, variableType, pos + l, nil
	}

	// read part - function(...) / (associative group...) / object.property / object.function()
	// conceptually, parenthesis grouping is a special case of anonymous identity function
	// NOTE: named expression types that contain a '.' are reserved for Object's only.
	if expr[pos] == '(' || (expr[pos] >= 'a' && expr[pos] <= 'z') || (expr[pos] >= 'A' && expr[pos] <= 'Z') {
		fname, lf, err := readNamedExpressionType(expr[pos:])
		switch {
		case errors.Is(err, errFunctionNameWithoutParens):
			if strings.Contains(fname, ".") {
				// user-define object property found.
				// note: we have already dealt with variableType above
				return fname, objectPropertyType, pos + lf, nil
			}
			// allow to continue so we can check alphanumerical operator names such as "And", "Or", etc
		case err != nil:
			return "", unknownType, 0, err
		default:
			fargs, la, err := readFunctionArguments(expr[pos+lf:])
			if err != nil {
				return "", unknownType, 0, err
			}
			if strings.Contains(fname, ".") {
				// user-defined object method found.
				return fname + fargs, objectMethodType, pos + lf + la, nil
			}
			return fname + fargs, functionType, pos + lf + la, nil
		}
	}

	// read part - object accessor (Dot operator)
	//
	// NOTE: object accessors are second degree to variables and functions
	// The allow to continue traversing an object by property or method.
	// The dot operator is used after any gal.entry that returns a value that can be treated as an object.
	// For example "Pi().Add(10).Sub(5)" is a valid expression because "Pi()" returns a gal.Value and
	// hence a Go object (be it struct or interface).
	if expr[pos] == '.' {
		// NOTE: we are keeping the leading dot in expr when submitting to readNamedExpressionType().
		// This means that `fname`` will always be of the form ".name" (i.e. with leading dot).
		// Remember that readNamedExpressionType is designed to read past 1 dot at most. In other words,
		// it will not read ".name.name2" but only ".name".
		fname, lf, err := readNamedExpressionType(expr[pos:])
		switch {
		case errors.Is(err, errFunctionNameWithoutParens):
			// property found on general purpose object
			return fname, objectAccessorByPropertyType, pos + lf, nil

		case err != nil:
			return "", unknownType, 0, err

		default:
			// method found on general purpose object
			fargs, la, err := readFunctionArguments(expr[pos+lf:])
			if err != nil {
				return "", unknownType, 0, err
			}
			return fname + fargs, objectAccessorByMethodType, pos + lf + la, nil
		}
	}

	// read part - operator
	if s, l := readOperator(expr[pos:]); l != 0 {
		if s == "+" || s == "-" {
			s, l = squashPlusMinusChain(expr[pos:]) // TODO: move this into readOperator()?
		}
		return s, operatorType, pos + l, nil
	}

	// read part - number
	// TODO: complex numbers are not supported - could be "native" or via function or perhaps even a specialised MultiValue?
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
		to++
		if expr[i] == '\\' {
			escapes += 1
			continue
		}
		if r == '"' && (escapes == 0 || escapes&1 == 0) {
			break
		}
		// TODO: perhaps we should collapse the `\`'s, here?

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
		to++
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

// The last "bool" return value of this function is an `ok` type bool.
// It is set to true if we successfull read a constant (i.e. "True" or "False", etc)
func readConstant(expr string) (string, int, exprType, bool) {
	to := 0

readString:
	for _, r := range expr {
		to++
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z':
			continue
		case isBlankSpace(r):
			// we read a potential constant name
			to-- // eject the space character we just read
			break readString
		default:
			// not a constant name
			return "", 0, unknownType, false
		}
	}

	switch expr[:to] {
	case "True", "False":
		// it's a Bool
		return expr[:to], to, boolType, true
	default:
		return "", 0, unknownType, false
	}
}

var errFunctionNameWithoutParens = errors.New("function without Parenthesis")

// 'name(...)' is a function call
// '()' is associative parenthesis grouping
func readNamedExpressionType(expr string) (string, int, error) {
	to := 0 // this could be an anonymous identity function (i.e. simple case of parenthesis grouping)

	dotCount := 0

	for _, r := range expr {
		// first: check for indication of function name (i.e. "name() or object.name()")
		if r == '(' {
			return expr[:to], to, nil
		}

		// second: check for indication of object property (i.e. "object.property.")
		if r == '.' {
			dotCount++
			if dotCount == 2 {
				break // this is an object property, no parenttheses
			}
		}

		if isBlankSpace(r) {
			break
		}
		to++
	}

	return expr[:to], to, errFunctionNameWithoutParens
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

		to++
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
		to++
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
	switch {
	case strings.HasPrefix(s, And.String()):
		return s[:3], 3

	case strings.HasPrefix(s, Power.String()),
		strings.HasPrefix(s, LShift.String()),
		strings.HasPrefix(s, RShift.String()),
		strings.HasPrefix(s, EqualTo.String()),
		strings.HasPrefix(s, NotEqualTo.String()),
		strings.HasPrefix(s, GreaterThanOrEqual.String()),
		strings.HasPrefix(s, LessThanOrEqual.String()),
		strings.HasPrefix(s, And2.String()),
		strings.HasPrefix(s, Or.String()),
		strings.HasPrefix(s, Or2.String()):
		return s[:2], 2

	case strings.HasPrefix(s, Plus.String()),
		strings.HasPrefix(s, Minus.String()),
		strings.HasPrefix(s, Divide.String()),
		strings.HasPrefix(s, Multiply.String()),
		strings.HasPrefix(s, Modulus.String()),
		strings.HasPrefix(s, GreaterThan.String()),
		strings.HasPrefix(s, LessThan.String()):
		return s[:1], 1

	default:
		return "", 0
	}
}

func isOperator(s string) bool {
	_, l := readOperator(s)
	return l != 0
}
