package gal

import "fmt"

// Undefined is a special gal.Value that indicates an undefined evaluation outcome.
//
// This can be as a first class citizen, when an error occurs
// (e.g. a '/' operator without the left hand side).
//
// All implementors of gal.Value also encapsulate an Undefined value.
// This ensures a default behaviour as defined by "Undefined"
// when none is available on the implementor.
// For instance, Bool does not support RShift() and does not implement it.
// However, since Bool encapsulates an Undefined value, it will return
// an Undefined value when RShift() is called on it.
type Undefined struct {
	reason string // optional
}

func NewUndefined() Undefined {
	return Undefined{}
}

func NewUndefinedWithReasonf(format string, a ...any) Undefined {
	return Undefined{
		reason: fmt.Sprintf(format, a...),
	}
}

func (Undefined) kind() entryKind {
	return unknownEntryKind
}

// Equal satisfies the external Equaler interface such as in testify assertions and the cmp package
func (u Undefined) Equal(other Undefined) bool {
	return u.reason == other.reason
}

func (u Undefined) EqualTo(other Value) Bool {
	return False
}

func (u Undefined) NotEqualTo(other Value) Bool {
	return True
}

func (u Undefined) GreaterThan(other Value) Bool {
	return False
}

func (u Undefined) GreaterThanOrEqual(other Value) Bool {
	return False
}

func (u Undefined) LessThan(other Value) Bool {
	return False
}

func (u Undefined) LessThanOrEqual(other Value) Bool {
	return False
}

func (u Undefined) Add(Value) Value {
	return u
}

func (u Undefined) Sub(Value) Value {
	return u
}

func (u Undefined) Multiply(Value) Value {
	return u
}

func (u Undefined) Divide(Value) Value {
	return u
}

func (u Undefined) PowerOf(Value) Value {
	return u
}

func (u Undefined) Mod(Value) Value {
	return u
}

func (u Undefined) LShift(Value) Value {
	return u
}

func (u Undefined) RShift(Value) Value {
	return u
}

func (Undefined) And(other Value) Bool {
	return Bool{Undefined: NewUndefinedWithReasonf("error: '%T/%s':'%s' cannot use And with Undefined", other, other.kind().String(), other.String())}
}

func (Undefined) Or(other Value) Bool {
	return Bool{Undefined: NewUndefinedWithReasonf("error: '%T/%s':'%s' cannot use Or with Undefined", other, other.kind().String(), other.String())}
}

func (u Undefined) String() string {
	if u.reason == "" {
		return "undefined"
	}
	return "undefined: " + u.reason
}

func (u Undefined) AsString() String {
	return NewString(u.String())
}

// The purpose of this method is to allow the user to check if a Value is undefined.
// For instance, if a Value is Number but the Undefined property's reason is not empty,
// it means that the Value is not a valid Number.
func (u Undefined) IsUndefined() bool {
	return u.reason == "" // NOTE: this is not quite accurate: an Undefined may not hold a reason
}
