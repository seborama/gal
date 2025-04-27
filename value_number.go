package gal

import (
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Number struct {
	Undefined
	value decimal.Decimal
}

func NewNumber(i int64, exp int32) Number {
	d := decimal.New(i, exp)

	return Number{value: d}
}

func NewNumberFromInt(i int64) Number {
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

// Equal satisfies the external Equaler interface such as in testify assertions and the cmp package
func (n Number) Equal(other Number) bool {
	return n.value.Equal(other.value)
}

func (n Number) Add(other Value) Value {
	if v, ok := other.(Numberer); ok {
		return Number{value: n.value.Add(v.Number().value)}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) Sub(other Value) Value {
	if v, ok := other.(Numberer); ok {
		return Number{
			value: n.value.Sub(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) Multiply(other Value) Value {
	if v, ok := other.(Numberer); ok {
		return Number{
			value: n.value.Mul(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) Divide(other Value) Value {
	if v, ok := other.(Numberer); ok {
		return Number{
			value: n.value.Div(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) PowerOf(other Value) Value {
	if v, ok := other.(Numberer); ok {
		return Number{
			value: n.value.Pow(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) Mod(other Value) Value {
	if v, ok := other.(Numberer); ok {
		return Number{
			value: n.value.Mod(v.Number().value),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) IntPart() Value {
	return Number{
		value: n.value.Truncate(0),
	}
}

func (n Number) LShift(other Value) Value {
	if v, ok := other.(Numberer); ok {
		if v.Number().value.IsNegative() {
			return NewUndefinedWithReasonf("invalid negative left shift")
		}
		if !v.Number().value.IsInteger() {
			return NewUndefinedWithReasonf("invalid non-integer left shift")
		}

		return Number{
			value: n.value.Mul(decimal.NewFromInt(2).Pow(v.Number().value)).Floor(),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

func (n Number) RShift(other Value) Value {
	if v, ok := other.(Numberer); ok {
		if v.Number().value.IsNegative() {
			return NewUndefinedWithReasonf("invalid negative right shift")
		}
		if !v.Number().value.IsInteger() {
			return NewUndefinedWithReasonf("invalid non-integer right shift")
		}

		return Number{
			value: n.value.Div(decimal.NewFromInt(2).Pow(v.Number().value)).Floor(),
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

func (n Number) Sqrt() Value {
	n, err := NewNumberFromString(
		new(big.Float).Sqrt(n.value.BigFloat()).String(),
	)
	if err != nil {
		return NewUndefinedWithReasonf("Sqrt:%s", err.Error())
	}

	return n
}

func (n Number) Tan() Number {
	return Number{
		value: n.value.Tan(),
	}
}

func (n Number) Ln(precision int32) Value {
	res, err := n.value.Ln(precision)
	if err != nil {
		return NewUndefinedWithReasonf("Ln:%s", err.Error())
	}

	return Number{
		value: res,
	}
}

func (n Number) Log(precision int32) Value {
	res, err := n.value.Ln(precision + 1)
	if err != nil {
		return NewUndefinedWithReasonf("Log:%s", err.Error())
	}

	res10, err := decimal.New(10, 0).Ln(precision + 1)
	if err != nil {
		return NewUndefinedWithReasonf("Log:%s", err.Error())
	}

	return Number{
		value: res.Div(res10).Truncate(precision),
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

func (n Number) Factorial() Value {
	if !n.value.IsInteger() || n.value.IsNegative() {
		return NewUndefinedWithReasonf("Factorial: requires a positive integer, cannot accept %s", n.String())
	}

	res := decimal.NewFromInt(1)

	one := decimal.NewFromInt(1)
	i := decimal.NewFromInt(2)
	for i.LessThanOrEqual(n.value) {
		res = res.Mul(i)
		i = i.Add(one)
	}

	return Number{
		value: res,
	}
}

func (n Number) LessThan(other Value) Bool {
	if v, ok := other.(Numberer); ok {
		return NewBool(n.value.LessThan(v.Number().value))
	}

	return False
}

func (n Number) LessThanOrEqual(other Value) Bool {
	if v, ok := other.(Numberer); ok {
		return NewBool(n.value.LessThanOrEqual(v.Number().value))
	}

	return False
}

func (n Number) EqualTo(other Value) Bool {
	if v, ok := other.(Numberer); ok {
		return NewBool(n.value.Equal(v.Number().value))
	}

	return False
}

func (n Number) NotEqualTo(other Value) Bool {
	return n.EqualTo(other).Not()
}

func (n Number) GreaterThan(other Value) Bool {
	if v, ok := other.(Numberer); ok {
		return NewBool(n.value.GreaterThan(v.Number().value))
	}

	return False
}

func (n Number) GreaterThanOrEqual(other Value) Bool {
	if v, ok := other.(Numberer); ok {
		return NewBool(n.value.GreaterThanOrEqual(v.Number().value))
	}

	return False
}

func (n Number) String() string {
	return n.value.String()
}

func (n Number) Bool() Bool {
	if n.value.IsZero() {
		return False
	}
	return True
}

func (n Number) AsString() String {
	return NewString(n.String())
}

func (n Number) Number() Number {
	return n
}

func (n Number) Float64() float64 {
	return n.value.InexactFloat64()
}

func (n Number) Int64() int64 {
	return n.value.IntPart()
}
