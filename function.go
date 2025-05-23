package gal

import (
	"fmt"
	"math"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
)

type FunctionalValue func(...Value) Value

func (fv FunctionalValue) String() string {
	return fmt.Sprintf("FunctionalValue @%p", fv)
}

type Function struct {
	Name     string
	Receiver Value // experimental concept: not used yet
	BodyFn   FunctionalValue
	Args     []Tree
}

func NewFunction(name string, bodyFn FunctionalValue, args ...Tree) Function {
	return Function{
		Name:   name,
		BodyFn: bodyFn,
		Args:   args,
	}
}

func (f Function) Calculate(val entry, op Operator, cfg *treeConfig) entry {
	if f.BodyFn == nil {
		// attempt to get body of a user-defined function
		// note: user-provided objects' methods are dealt with by ObjectMethod.Calculate
		f.BodyFn = cfg.Function(f.Name)
	}

	rhsVal := f.Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables), WithObjects(cfg.objects))
	if u, ok := rhsVal.(Undefined); ok {
		return u
	}

	if val == nil {
		return rhsVal
	}

	//nolint:errcheck // life's too short to check for type assertion success here
	val = calculate(val.(Value), op, rhsVal)

	return val
}

func (f Function) String() string {
	args := lo.Map(f.Args, func(item Tree, index int) string {
		return strings.TrimRight(item.String(), "\n")
	})
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
}

// Equal satisfies the external Equaler interface such as in testify assertions and the cmp package
func (f Function) Equal(other Function) bool {
	return f.Name == other.Name &&
		f.BodyFn.String() == other.BodyFn.String() &&
		cmp.Equal(f.Args, other.Args)
}

func (f Function) Eval(opts ...treeOption) Value {
	var args []Value

	// TODO: how about adding a Receiver property to Function?
	// ...   It would be populated with the user-defined object or runtime "LHS" accessor, when one is present.
	// ...   This method would be responsible to populate BodyFn by calling ObjectGetMethod on the receiver.
	if f.Receiver != nil {
		bodyFn, ok := ObjectGetMethod(f.Receiver, f.Name)
		if !ok {
			return NewUndefinedWithReasonf("unknown method '%s' for receiver '%T'", f.Name, f.Receiver)
		}
		f.BodyFn = bodyFn
	}

	for _, a := range f.Args {
		args = append(args, a.Eval(opts...))
	}

	if f.BodyFn == nil {
		return NewUndefinedWithReasonf("unknown function '%s'", f.Name)
	}

	return f.BodyFn(args...)
}

var builtInFunctions = map[string]FunctionalValue{
	"pi":        Pi,
	"factorial": Factorial,
	"cos":       Cos,
	"sin":       Sin,
	"tan":       Tan,
	"sqrt":      Sqrt,
	"floor":     Floor,
	"trunc":     Trunc,
	"ln":        Ln,
	"log":       Log,
	"eval":      Eval,
}

// BuiltInFunction returns a built-in function body if known.
// It returns `nil` when no built-in function exists by the specified name.
// This signals the Evaluator to attempt to find a user defined function.
func BuiltInFunction(name string) FunctionalValue {
	if builtInFunctions == nil {
		return nil
	}

	// note: for now function names are arbitrarily case-insensitive
	lowerName := strings.ToLower(name)

	bodyFn, ok := builtInFunctions[lowerName]
	if ok {
		return bodyFn
	}

	return nil
}

// Pi returns the Value of math.Pi.
// TODO: this could likely be turned into a constant.
func Pi(args ...Value) Value {
	if len(args) != 0 {
		return NewUndefinedWithReasonf("pi() requires no argument, got %d", len(args))
	}

	return NewNumberFromFloat(math.Pi)
}

// PiLong returns a value of Pi with many more digits than Pi.
// TODO: this could likely be turned into a constant.
func PiLong(args ...Value) Value {
	if len(args) != 0 {
		return NewUndefinedWithReasonf("pi() requires no argument, got %d", len(args))
	}

	pi, _ := NewNumberFromString(Pi51199) //nolint:errcheck

	return pi
}

// Factorial returns the factorial of the provided argument.
func Factorial(args ...Value) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("factorial() requires 1 argument, got %d", len(args))
	}
	if v, ok := args[0].(Numberer); ok {
		return v.Number().Factorial()
	}

	return NewUndefinedWithReasonf("factorial(): invalid argument type '%s'", args[0].String())
}

// Cos returns the cosine.
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

// Sin returns the sine.
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

// Tan returns the tangent.
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

// Ln returns the natural logarithm of d.
func Ln(args ...Value) Value {
	if len(args) != 2 {
		return NewUndefinedWithReasonf("ln() requires 2 arguments, got %d", len(args))
	}

	if v, ok := args[0].(Numberer); ok {
		if p, ok := args[1].(Numberer); ok {
			return v.Number().Ln(int32(p.Number().value.IntPart())) //nolint:gosec
		}
	}

	return NewUndefinedWithReasonf("ln(): invalid argument type '%s'", args[0].String())
}

// Log returns the logarithm base 10 of d.
func Log(args ...Value) Value {
	if len(args) != 2 {
		return NewUndefinedWithReasonf("log() requires 2 arguments, got %d", len(args))
	}

	if v, ok := args[0].(Numberer); ok {
		if p, ok := args[1].(Numberer); ok {
			return v.Number().Log(int32(p.Number().value.IntPart())) //nolint:gosec // ignoring overflow conversion
		}
	}

	return NewUndefinedWithReasonf("log(): invalid argument type '%s'", args[0].String())
}

// Sqrt returns the square root.
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

// return the floor.
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
		return v.Number().Trunc(int32(precision.Number().value.IntPart())) //nolint:gosec // ignoring overflow conversion
	}

	return NewUndefinedWithReasonf("trunc(): invalid argument #1 '%s'", argVal.String())
}

func Eval(args ...Value) Value {
	if len(args) != 1 {
		return NewUndefinedWithReasonf("eval() requires 1 argument1, got %d: '%v'", len(args), args)
	}

	argVal := args[0]

	if v, ok := argVal.(Evaler); ok {
		return v.Eval()
	}

	return argVal
}
