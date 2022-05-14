package gal

import (
	"fmt"
	"math"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type FunctionalValue func(...Tree) Value

func (fv FunctionalValue) String() string {
	// TODO: This is only to support the tests (see Function.Equal)).
	//       This is not elegant but so far the only solution I have to compare
	//       Functions in the tests.
	return fmt.Sprintf("FunctionalValue @%p", fv)
}

type Function struct {
	BodyFn FunctionalValue
	Args   []Tree
}

func NewFunction(bodyFn FunctionalValue, args ...Tree) Function {
	return Function{
		BodyFn: bodyFn,
		Args:   args,
	}
}

func (Function) kind() entryKind {
	return functionEntryKind
}

func (f Function) Equal(other Function) bool {
	// TODO: This is only to support the tests.
	//       This is not elegant but so far the only solution I have to compare
	//       Functions in the tests.
	return f.BodyFn.String() == other.BodyFn.String() &&
		cmp.Equal(f.Args, other.Args)
}

func (f Function) Eval() Value {
	return f.BodyFn(f.Args...)
}

var preDefinedFunctions = map[string]FunctionalValue{
	"pi":    Pi,
	"cos":   Cos,
	"sin":   Sin,
	"tan":   Tan,
	"sqrt":  Sqrt,
	"floor": Floor,
	"trunc": Trunc,
}

func PreDefinedFunction(name string, userFunctions Functions) FunctionalValue {
	// note: for now function names are arbitrarily case-insensitive
	lowerName := strings.ToLower(name)

	bodyFn, ok := preDefinedFunctions[lowerName]
	if ok {
		return bodyFn
	}

	bodyFn, ok = userFunctions.Function(lowerName)
	if ok {
		return bodyFn
	}

	return func(...Tree) Value {
		return NewUndefinedWithReasonf("unknown function '%s'", name)
	}
}

func Pi(args ...Tree) Value {
	if len(args) != 0 {
		return NewUndefinedWithReasonf("pi() requires no argument, got %d", len(args))
	}

	return NewNumberFromFloat(math.Pi)
}

func Cos(args ...Tree) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("cos() requires 1 argument, got %d", len(args))
	}

	argVal := args[0].Eval()
	if v, ok := argVal.(numberer); ok {
		return v.Number().Cos()
	}

	return NewUndefinedWithReasonf("cos(): invalid argument type '%T'", args[0])
}

func Sin(args ...Tree) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("sin() requires 1 argument, got %d", len(args))
	}

	argVal := args[0].Eval()
	if v, ok := argVal.(numberer); ok {
		return v.Number().Sin()
	}

	return NewUndefinedWithReasonf("sin(): invalid argument type '%T'", args[0])
}

func Tan(args ...Tree) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("tan() requires 1 argument, got %d", len(args))
	}

	argVal := args[0].Eval()
	if v, ok := argVal.(numberer); ok {
		return v.Number().Tan()
	}

	return NewUndefinedWithReasonf("tan(): invalid argument type '%T'", args[0])
}

func Sqrt(args ...Tree) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("sqrt() requires 1 argument, got %d", len(args))
	}

	argVal := args[0].Eval()
	if v, ok := argVal.(numberer); ok {
		return v.Number().Sqrt()
	}

	return NewUndefinedWithReasonf("sqrt(): invalid argument type '%T'", args[0])
}

func Floor(args ...Tree) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("floor() requires 1 argument, got %d", len(args))
	}

	argVal := args[0].Eval()
	if v, ok := argVal.(numberer); ok {
		return v.Number().Floor()
	}

	return NewUndefinedWithReasonf("floor(): invalid argument type '%T'", args[0])
}

func Trunc(args ...Tree) Value {
	if len(args) != 2 {
		return NewUndefinedWithReasonf("trunc() requires 2 arguments, got %d: '%v'", len(args), args)
	}

	argVal := args[0].Eval()

	argPrecision := args[1].Eval()
	precision, ok := argPrecision.(numberer)
	if !ok {
		return NewUndefinedWithReasonf("trunc() requires precision (argument #1) to be a number, got %s", argPrecision.String())
	}

	if v, ok := argVal.(numberer); ok {
		return v.Number().Trunc(int32(precision.Number().value.IntPart()))
	}

	return NewUndefinedWithReasonf("trunc(): invalid argument #2 '%s'", argVal.String())
}
