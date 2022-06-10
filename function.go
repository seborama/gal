package gal

import (
	"fmt"
	"math"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type FunctionalValue func(...Value) Value

func (fv FunctionalValue) String() string {
	// TODO: This is only to support the tests (see Function.Equal)).
	//       This is not elegant but so far the only solution I have to compare
	//       Functions in the tests.
	return fmt.Sprintf("FunctionalValue @%p", fv)
}

type Function struct {
	Name   string
	BodyFn FunctionalValue
	Args   []Tree
}

func NewFunction(name string, bodyFn FunctionalValue, args ...Tree) Function {
	return Function{
		Name:   name,
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
	return f.Name == other.Name &&
		f.BodyFn.String() == other.BodyFn.String() &&
		cmp.Equal(f.Args, other.Args)
}

func (f Function) Eval(opts ...treeOption) Value {
	var args []Value

	for _, a := range f.Args {
		args = append(args, a.Eval(opts...))
	}

	if f.BodyFn == nil {
		return NewUndefinedWithReasonf("unknown function '%s'", f.Name)
	}

	return f.BodyFn(args...)
}

var preDefinedFunctions = map[string]FunctionalValue{
	"pi":        Pi,
	"factorial": Factorial,
	"cos":       Cos,
	"sin":       Sin,
	"tan":       Tan,
	"sqrt":      Sqrt,
	"floor":     Floor,
	"trunc":     Trunc,
}

// PreDefinedFunction returns a pre-defined function body if known.
// It returns `nil` when no pre-defined function exists by the specified name.
// This signals the Evaluator to attempt to find a user defined function.
func PreDefinedFunction(name string) FunctionalValue {
	// note: for now function names are arbitrarily case-insensitive
	lowerName := strings.ToLower(name)

	bodyFn, ok := preDefinedFunctions[lowerName]
	if ok {
		return bodyFn
	}

	return nil
}

func UserDefinedFunction(name string, userFunctions Functions) FunctionalValue {
	// note: for now function names are arbitrarily case-insensitive
	lowerName := strings.ToLower(name)

	return userFunctions.Function(lowerName)
}

func Pi(args ...Value) Value {
	if len(args) != 0 {
		return NewUndefinedWithReasonf("pi() requires no argument, got %d", len(args))
	}

	return NewNumberFromFloat(math.Pi)
}

func Factorial(args ...Value) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("factorial() requires 1 argument, got %d", len(args))
	}
	if v, ok := args[0].(Numberer); ok {
		return v.Number().Factorial()
	}

	return NewUndefinedWithReasonf("factorial(): invalid argument type '%s'", args[0].String())
}

func Cos(args ...Value) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("cos() requires 1 argument, got %d", len(args))
	}

	argVal := args[0]
	if v, ok := argVal.(Numberer); ok {
		return v.Number().Cos()
	}

	return NewUndefinedWithReasonf("cos(): invalid argument type '%s'", args[0].String())
}

func Sin(args ...Value) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("sin() requires 1 argument, got %d", len(args))
	}

	argVal := args[0]
	if v, ok := argVal.(Numberer); ok {
		return v.Number().Sin()
	}

	return NewUndefinedWithReasonf("sin(): invalid argument type '%s'", args[0].String())
}

func Tan(args ...Value) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("tan() requires 1 argument, got %d", len(args))
	}

	argVal := args[0]
	if v, ok := argVal.(Numberer); ok {
		return v.Number().Tan()
	}

	return NewUndefinedWithReasonf("tan(): invalid argument type '%s'", args[0].String())
}

func Sqrt(args ...Value) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("sqrt() requires 1 argument, got %d", len(args))
	}

	argVal := args[0]
	if v, ok := argVal.(Numberer); ok {
		return v.Number().Sqrt()
	}

	return NewUndefinedWithReasonf("sqrt(): invalid argument type '%T'", args[0])
}

func Floor(args ...Value) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("floor() requires 1 argument, got %d", len(args))
	}

	argVal := args[0]
	if v, ok := argVal.(Numberer); ok {
		return v.Number().Floor()
	}

	return NewUndefinedWithReasonf("floor(): invalid argument type '%T'", args[0])
}

func Trunc(args ...Value) Value {
	if len(args) != 2 {
		return NewUndefinedWithReasonf("trunc() requires 2 arguments, got %d: '%v'", len(args), args)
	}

	argVal := args[0]

	argPrecision := args[1]
	precision, ok := argPrecision.(Numberer)
	if !ok {
		return NewUndefinedWithReasonf("trunc() requires precision (argument #2) to be a number, got %s", argPrecision.String())
	}

	if v, ok := argVal.(Numberer); ok {
		return v.Number().Trunc(int32(precision.Number().value.IntPart()))
	}

	return NewUndefinedWithReasonf("trunc(): invalid argument #1 '%s'", argVal.String())
}
