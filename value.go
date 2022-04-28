package gal

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type String struct {
	value string
}

func NewString(s string) String {
	return String{value: s}
}

func (String) kind() entryKind {
	return valueEntryKind
}

func (s String) Equal(other String) bool {
	return s.value == other.value
}

func (s String) Add(other Value) Value {
	switch v := other.(type) {
	case String:
		return String{value: s.value + v.value}
	}

	v, ok := other.(stringer)
	if !ok {
		return Undefined{}
	}

	return String{value: s.value + v.String()}
}

func (s String) Times(other Value) Value {
	switch v := other.(type) {
	case Number:
		if !v.value.IsInteger() {
			return Undefined{}
		}

		return String{value: strings.Repeat(s.value, int(v.value.IntPart()))}
	}

	v, ok := other.(stringer)
	if !ok {
		return Undefined{}
	}

	return String{value: s.value + v.String()}
}

func (s String) String() string {
	return s.value
}

type Number struct {
	value decimal.Decimal
}

func NewNumber(i int64) Number {
	d := decimal.NewFromInt(i)

	return Number{value: d}
}

func NewNumberFromFloat(f float64) Number {
	d := decimal.NewFromFloat(f)

	return Number{value: d}
}

func NewNumberFromString(s string) (Number, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Number{}, errors.WithStack(err)
	}

	return Number{value: d}, nil
}

func (Number) kind() entryKind {
	return valueEntryKind
}

func (n Number) Equal(other Number) bool {
	return n.value.Equal(other.value)
}

func (n Number) Add(other Value) Value {
	switch v := other.(type) {
	case Number:
		return Number{value: n.value.Add(v.value)}
	}

	v, ok := other.(numberer)
	if !ok {
		return Undefined{}
	}

	return Number{
		value: n.value.Add(v.Number()),
	}
}

func (n Number) Times(other Value) Value {
	switch v := other.(type) {
	case Number:
		return Number{value: n.value.Mul(v.value)}
	}

	v, ok := other.(numberer)
	if !ok {
		return Undefined{}
	}

	return Number{
		value: n.value.Mul(v.Number()),
	}
}

func (n Number) String() string {
	return n.value.String()
}

type Undefined struct {
	reason string // optional
}

func NewUndefined() Undefined {
	return Undefined{}
}

func NewUndefinedWithReason(reason string) Undefined {
	return Undefined{
		reason: reason,
	}
}

func (u Undefined) Equal(other Undefined) bool {
	return u.reason == other.reason
}

func (Undefined) Add(Value) Value {
	return Undefined{}
}

func (Undefined) Times(Value) Value {
	return Undefined{}
}

func (Undefined) String() string {
	return "undefined"
}