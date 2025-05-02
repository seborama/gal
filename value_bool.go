package gal

import "github.com/pkg/errors"

type Bool struct {
	Undefined
	value bool
}

func NewBool(b bool) Bool {
	return Bool{value: b}
}

// NOTE: another option would be to return:
// Bool{Undefined: NewUndefinedWithReasonf("cannot convert '%s' to Bool", s)}
func NewBoolFromString(s string) (Bool, error) {
	switch s {
	case "True":
		return True, nil
	case "False":
		return False, nil
	default:
		return False, errors.Errorf("'%s' cannot be converted to a Bool", s)
	}
}

// Equal satisfies the external Equaler interface such as in testify assertions and the cmp package
func (b Bool) Equal(other Bool) bool {
	return b.value == other.value
}

func (b Bool) EqualTo(other Value) Bool {
	if v, ok := other.(Booler); ok {
		return NewBool(b.value == v.Bool().value)
	}
	return False
}

func (b Bool) NotEqualTo(other Value) Bool {
	return b.EqualTo(other).Not()
}

func (b Bool) Not() Bool {
	return NewBool(!b.value)
}

func (b Bool) And(other Value) Bool {
	if v, ok := other.(Booler); ok {
		return NewBool(b.value && v.Bool().value)
	}
	return False // NOTE: should Bool be a Maybe?
}

func (b Bool) Or(other Value) Bool {
	if v, ok := other.(Booler); ok {
		return NewBool(b.value || v.Bool().value)
	}
	return False // NOTE: should Bool be a Maybe?
}

func (b Bool) Bool() Bool {
	return b
}

func (b Bool) String() string {
	if b.value {
		return "True"
	}
	return "False"
}

func (b Bool) Number() Number {
	if b.value {
		return NewNumberFromInt(1)
	}
	return NewNumberFromInt(0)
}

func (b Bool) AsString() String {
	return NewString(b.String())
}

var (
	False = NewBool(false)
	True  = NewBool(true)
)
