package gal

import "fmt"

type Value interface {
	valueCalculation
	valueComparison
	valueLogic
	valueHelper
	undefinedChecker
}

type valueCalculation interface {
	Add(Value) Value
	Sub(Value) Value
	Multiply(Value) Value
	Divide(Value) Value
	PowerOf(Value) Value
	Mod(Value) Value
	LShift(Value) Value
	RShift(Value) Value
}

type valueComparison interface {
	LessThan(Value) Bool
	LessThanOrEqual(Value) Bool
	EqualTo(Value) Bool
	NotEqualTo(Value) Bool
	GreaterThan(Value) Bool
	GreaterThanOrEqual(Value) Bool
}

type valueLogic interface {
	And(Value) Bool
	Or(Value) Bool
}

type valueHelper interface {
	Stringer
	fmt.Stringer
	entry
}

type undefinedChecker interface {
	// TODO: IsUndefined somewhat mimics a "Maybe" monad in functional programming:
	// ...   e.g. if a Bool has its Undefined value set, IsUndefined will return true.
	// ...   Instead of using the Bool, we should unwrap the Undefined and use it: this is not
	// ...   implemented yet!
	IsUndefined() bool
}

type Stringer interface {
	AsString() String // name is not String so to not clash with fmt.Stringer interface
}

type Numberer interface {
	Number() Number
}

type Booler interface {
	Bool() Bool
}

type Evaler interface {
	Eval() Value
}

func ToValue(value any) Value {
	v, _ := toValue(value) // ignore "ok" because we are sure it is a valid Value, be it Undefined or not.
	return v
}

func toValue(value any) (Value, bool) {
	v, err := goAnyToGalType(value)
	if err != nil {
		return NewUndefinedWithReasonf("value type %T - %s", value, err.Error()), false
	}
	return v, true
}

func ToNumber(val Value) Number {
	n, ok := val.(Numberer)
	if !ok {
		return Number{Undefined: NewUndefinedWithReasonf("value type %T - cannot convert to Number", val)}
	}
	return n.Number()
}

func ToString(val Value) String {
	return val.AsString()
}

func ToBool(val Value) Bool {
	b, ok := val.(Booler)
	if !ok {
		return Bool{Undefined: NewUndefinedWithReasonf("value type %T - cannot convert to Bool", val)}
	}
	return b.Bool()
}
