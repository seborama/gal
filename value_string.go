package gal

import "strings"

type String struct {
	Undefined
	value string
}

func NewString(s string) String {
	return String{value: s}
}

func (String) kind() entryKind {
	return valueEntryKind
}

// Equal satisfies the external Equaler interface such as in testify assertions and the cmp package
func (s String) Equal(other String) bool {
	return s.value == other.value
}

func (s String) LessThan(other Value) Bool {
	if v, ok := other.(Stringer); ok {
		return NewBool(s.value < v.AsString().value)
	}

	return False
}

func (s String) LessThanOrEqual(other Value) Bool {
	if v, ok := other.(Stringer); ok {
		return NewBool(s.value <= v.AsString().value)
	}

	return False
}

func (s String) EqualTo(other Value) Bool {
	if v, ok := other.(Stringer); ok {
		return NewBool(s.value == v.AsString().value) // beware to compare what's comparable: do NOT use s.value == v.String() because String() may decorate the value (see String and MultiValue for example)
	}

	return False
}

func (s String) NotEqualTo(other Value) Bool {
	return s.EqualTo(other).Not()
}

func (s String) GreaterThan(other Value) Bool {
	if v, ok := other.(Stringer); ok {
		return NewBool(s.value > v.AsString().value)
	}

	return False
}

func (s String) GreaterThanOrEqual(other Value) Bool {
	if v, ok := other.(Stringer); ok {
		return NewBool(s.value >= v.AsString().value)
	}

	return False
}

func (s String) Add(other Value) Value {
	if v, ok := other.(Stringer); ok {
		return String{value: s.value + v.AsString().value}
	}

	return NewUndefinedWithReasonf("cannot Add non-string to a string")
}

func (s String) Multiply(other Value) Value {
	if v, ok := other.(Numberer); ok {
		count := v.Number().value
		if !count.IsInteger() || count.IsNegative() {
			return NewUndefinedWithReasonf("String.Multiply: invalid repeat count: %s", count.String())
		}
		n := count.IntPart()
		if int64(int(n)) != n { // overflow check
			return NewUndefinedWithReasonf("String.Multiply: repeat count overflows on this architecture")
		}
		return String{value: strings.Repeat(s.value, int(n))}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

// TODO: add test to confirm this is correct!
func (s String) LShift(other Value) Value {
	v, ok := other.(Numberer)
	if !ok {
		return NewUndefinedWithReasonf("NaN: %s", other.String())
	}

	if v.Number().value.IsNegative() {
		return NewUndefinedWithReasonf("invalid negative left shift")
	}
	if !v.Number().value.IsInteger() {
		return NewUndefinedWithReasonf("invalid non-integer left shift")
	}

	idx64 := v.Number().value.IntPart()
	if idx64 < 0 {
		return NewUndefinedWithReasonf("left shift [%s]: out of range", other.String())
	}
	if idx64 > int64(len(s.value)) {
		return String{}
	}

	return String{value: s.value[int(idx64):]}
}

// TODO: add test to confirm this is correct!
func (s String) RShift(other Value) Value {
	v, ok := other.(Numberer)
	if !ok {
		return NewUndefinedWithReasonf("NaN: %s", other.String())
	}

	if v.Number().value.IsNegative() {
		return NewUndefinedWithReasonf("invalid negative right shift")
	}
	if !v.Number().value.IsInteger() {
		return NewUndefinedWithReasonf("invalid non-integer right shift")
	}

	shift := v.Number().value.IntPart()
	if shift < 0 {
		return NewUndefinedWithReasonf("right shift [%s]: out of range", other.String())
	}
	limit := int64(len(s.value))
	if shift > limit {
		return String{}
	}

	return String{value: s.value[:int64(len(s.value))-v.Number().value.IntPart()]}
}

func (s String) String() string {
	return `"` + s.value + `"`
}

func (s String) RawString() string {
	return s.value
}

func (s String) AsString() String {
	return s
}

func (s String) Number() Number {
	n, err := NewNumberFromString(s.value) // beware that `.String()` may decorate the value!!
	if err != nil {
		return Number{Undefined: NewUndefinedWithReasonf("cannot convert %s to Number: %s", s.String(), err.Error())}
	}

	return n
}

func (s String) Eval() Value {
	tree, err := NewTreeBuilder().FromExpr(s.value)
	if err != nil {
		return s
	}

	return tree.Eval()
}
