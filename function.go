package gal

import (
	"math"
	"strings"
)

type Function struct {
	Name string
	Args []Tree
}

func NewFunction(name string, args ...Tree) Function {
	return Function{
		Name: name,
		Args: args,
	}
}

func (Function) kind() entryKind {
	return functionEntryKind
}

func (f Function) String() string {
	return string(f.Name)
}

func (f Function) argsLen() int {
	return len(f.Args)
}

func (f Function) Eval() Value {
	// note: for now function names are arbitrarily case-insensitive
	switch strings.ToLower(f.Name) {
	case "pi":
		if f.argsLen() != 0 {
			return NewUndefinedWithReasonf("%s requires no parameter, got %d", f.Name, f.argsLen())
		}
		return NewNumberFromFloat(math.Pi)

	case "cos":
		if f.argsLen() != 1 {
			return NewUndefinedWithReasonf("%s requires 1 parameter, got %d", f.Name, f.argsLen())
		}
		argVal := f.Args[0].Eval()
		if v, ok := argVal.(numberer); ok {
			return v.Number().Cos()
		}

	case "floor":
		if f.argsLen() != 1 {
			return NewUndefinedWithReasonf("%s requires 1 parameter, got %d", f.Name, f.argsLen())
		}
		argVal := f.Args[0].Eval()
		if v, ok := argVal.(numberer); ok {
			return v.Number().Floor()
		}

	case "sin":
		if f.argsLen() != 1 {
			return NewUndefinedWithReasonf("%s requires 1 parameter, got %d", f.Name, f.argsLen())
		}
		argVal := f.Args[0].Eval()
		if v, ok := argVal.(numberer); ok {
			return v.Number().Sin()
		}

	case "sqrt":
		if f.argsLen() != 1 {
			return NewUndefinedWithReasonf("%s requires 1 parameter, got %d", f.Name, f.argsLen())
		}
		argVal := f.Args[0].Eval()
		if v, ok := argVal.(numberer); ok {
			return v.Number().Sqrt()
		}

	case "tan":
		if f.argsLen() != 1 {
			return NewUndefinedWithReasonf("%s requires 1 parameter, got %d", f.Name, f.argsLen())
		}
		argVal := f.Args[0].Eval()
		if v, ok := argVal.(numberer); ok {
			return v.Number().Tan()
		}

	case "trunc":
		if f.argsLen() != 2 {
			return NewUndefinedWithReasonf("%s requires 2 parameters, got %d", f.Name, f.argsLen())
		}
		argPrecision := f.Args[0].Eval()
		precision, ok := argPrecision.(numberer)
		if !ok {
			return NewUndefinedWithReasonf("%s requires precision (argument #1) to be a number, got %s", f.Name, argPrecision.String())
		}
		argVal := f.Args[1].Eval()
		if v, ok := argVal.(numberer); ok {
			return v.Number().Trunc(int32(precision.Number().value.IntPart()))
		}
		return NewUndefinedWithReasonf("%s requires value (argument #2) to be a number, got %s", f.Name, argVal.String())
	}

	return NewUndefinedWithReasonf("unknown function '%s'", f.Name)
}
