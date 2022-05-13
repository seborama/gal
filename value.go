package gal

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type stringer interface {
	String() string
}

type numberer interface {
	Number() Number
}

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
	if v, ok := other.(stringer); ok {
		return String{value: s.value + v.String()}
	}

	return NewUndefinedWithReasonf("cannot Add non-string to a string")
}

func (s String) Sub(other Value) Value {
	return NewUndefinedWithReasonf("cannot Sub from string")
}

func (s String) Multiply(other Value) Value {
	if v, ok := other.(numberer); ok {
		return String{
			value: strings.Repeat(s.value, int(v.Number().value.IntPart())),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (s String) Divide(other Value) Value {
	return Undefined{}
}

func (s String) PowerOf(Value) Value {
	return Undefined{}
}

func (s String) Mod(Value) Value {
	return Undefined{}
}

func (s String) String() string {
	return s.value
}

func (s String) Number() Number {
	n, err := NewNumberFromString(s.String())
	if err != nil {
		panic(err) // TODO :-/
	}

	return n
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
	if v, ok := other.(numberer); ok {
		return Number{value: n.value.Add(v.Number().value)}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) Sub(other Value) Value {
	switch v := other.(type) {
	case Number:
		return Number{value: n.value.Sub(v.value)}
	}

	if v, ok := other.(numberer); ok {
		return Number{
			value: n.value.Sub(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) Multiply(other Value) Value {
	if v, ok := other.(numberer); ok {
		return Number{
			value: n.value.Mul(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) Divide(other Value) Value {
	if v, ok := other.(numberer); ok {
		return Number{
			value: n.value.Div(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) PowerOf(other Value) Value {
	if v, ok := other.(numberer); ok {
		return Number{
			value: n.value.Pow(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) Mod(other Value) Value {
	if v, ok := other.(numberer); ok {
		return Number{
			value: n.value.Mod(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) Neg() Number {
	return Number{
		value: n.value.Neg(),
	}
}

func (n Number) Sin() Number {
	return Number{
		value: n.value.Sin(),
	}
}

func (n Number) Cos() Number {
	return Number{
		value: n.value.Cos(),
	}
}

func (n Number) Sqrt() Number {
	n, err := NewNumberFromString(
		new(big.Float).Sqrt(n.value.BigFloat()).String(),
	)
	if err != nil {
		panic(err) // TODO: :-/
	}

	return n
}

func (n Number) Tan() Number {
	return Number{
		value: n.value.Tan(),
	}
}

func (n Number) Floor() Number {
	return Number{
		value: n.value.Floor(),
	}
}

func (n Number) Trunc(precision int32) Number {
	return Number{
		value: n.value.Truncate(precision),
	}
}

func (n Number) String() string {
	return n.value.String()
}

func (n Number) Number() Number {
	return n
}

type Undefined struct {
	reason string // optional
}

func NewUndefined() Undefined {
	return Undefined{}
}

func NewUndefinedWithReasonf(format string, a ...interface{}) Undefined {
	return Undefined{
		reason: fmt.Sprintf(format, a...),
	}
}

func (Undefined) kind() entryKind {
	return unknownEntryKind
}

func (u Undefined) Equal(other Undefined) bool {
	return u.reason == other.reason
}

func (Undefined) Add(Value) Value {
	return Undefined{}
}

func (Undefined) Sub(Value) Value {
	return Undefined{}
}

func (Undefined) Multiply(Value) Value {
	return Undefined{}
}

func (Undefined) Divide(Value) Value {
	return Undefined{}
}

func (Undefined) PowerOf(Value) Value {
	return Undefined{}
}

func (Undefined) Mod(Value) Value {
	return Undefined{}
}

func (u Undefined) String() string {
	if u.reason == "" {
		return "undefined"
	}
	return "undefined: " + u.reason
}
