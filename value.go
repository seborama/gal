package gal

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

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
	v, _ := toValue(value)
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
	//nolint:errcheck // life's too short to check for type assertion success here
	return val.(Numberer).Number() // may panic
}

func ToString(val Value) String {
	return val.AsString()
}

func ToBool(val Value) Bool {
	//nolint:errcheck // life's too short to check for type assertion success here
	return val.(Booler).Bool() // may panic
}

// MultiValue is a container of zero or more Value's.
// For the time being, it is only usable and useful with functions.
// Functions can accept a MultiValue, and also return a MultiValue.
// This allows a function to effectively return multiple values as a MultiValue.
// TODO: we could add a syntax to instantiate a MultiValue within an expression.
// ...   perhaps along the lines of [[v1 v2 ...]] or simply a built-in function such as
// ...   MultiValue(...) - nothing stops the user from creating their own for now :-)
//
// TODO: implement other methods such as Add, LessThan, etc (if meaningful)
type MultiValue struct {
	Undefined
	values []Value
}

func NewMultiValue(values ...Value) MultiValue {
	return MultiValue{values: values}
}

func (MultiValue) kind() entryKind {
	return valueEntryKind
}

// Equal satisfies the external Equaler interface such as in testify assertions and the cmp package
// Note that the current implementation defines equality as values matching and in order they appear.
func (m MultiValue) Equal(other MultiValue) bool {
	if m.Size() != other.Size() {
		return false
	}

	for i := range m.values {
		// TODO: add test to confirm this is correct!
		if m.values[i].NotEqualTo(other.values[i]) == False {
			return false
		}
	}

	return true
}

func (m MultiValue) String() string {
	var vals []string
	for _, val := range m.values {
		vals = append(vals, val.String())
	}
	return strings.Join(vals, `,`)
}

func (m MultiValue) AsString() String {
	return NewString(m.String())
}

func (m MultiValue) Get(i int) Value {
	if i > len(m.values) {
		return NewUndefinedWithReasonf("out of bounds: trying to get arg #%d on MultiValue that has %d arguments", i, len(m.values))
	}

	return m.values[i]
}

func (m MultiValue) Size() int {
	return len(m.values)
}

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
		return String{
			value: strings.Repeat(s.value, int(v.Number().value.IntPart())),
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

// TODO: add test to confirm this is correct!
func (s String) LShift(other Value) Value {
	if v, ok := other.(Numberer); ok {
		if v.Number().value.IsNegative() {
			return NewUndefinedWithReasonf("invalid negative left shift")
		}
		if !v.Number().value.IsInteger() {
			return NewUndefinedWithReasonf("invalid non-integer left shift")
		}

		return String{
			value: s.value[v.Number().value.IntPart():],
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
}

// TODO: add test to confirm this is correct!
func (s String) RShift(other Value) Value {
	if v, ok := other.(Numberer); ok {
		if v.Number().value.IsNegative() {
			return NewUndefinedWithReasonf("invalid negative left shift")
		}
		if !v.Number().value.IsInteger() {
			return NewUndefinedWithReasonf("invalid non-integer left shift")
		}

		return String{
			value: s.value[:int64(len(s.value))-v.Number().value.IntPart()],
		}
	}

	return NewUndefinedWithReasonf("NaN: %s", other.String())
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
		panic(err) // TODO :-/
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

type Bool struct {
	Undefined
	value bool
}

func NewBool(b bool) Bool {
	return Bool{value: b}
}

// TODO: another option would be to return a Value and hence allow Undefined when neither True nor False is provided.
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

func (Bool) kind() entryKind {
	return valueEntryKind
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
	return False // TODO: should Bool be a Maybe?
}

func (b Bool) Or(other Value) Bool {
	if v, ok := other.(Booler); ok {
		return NewBool(b.value || v.Bool().value)
	}
	return False // TODO: should Bool be a Maybe?
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

func (Undefined) LShift(Value) Value {
	return Undefined{}
}

func (Undefined) RShift(Value) Value {
	return Undefined{}
}

func (Undefined) And(other Value) Bool {
	// perhaps this should be a panic... Or else Bool should be a Maybe?
	return False
}

func (Undefined) Or(other Value) Bool {
	// perhaps this should be a panic... Or else Bool should be a Maybe?
	return False
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
