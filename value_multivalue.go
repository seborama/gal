package gal

import "strings"

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

// Equal satisfies the external Equaler interface such as in `testify` assertions and the `cmp` package
// Note that the current implementation defines equality as values matching and in order they appear.
func (m MultiValue) Equal(other MultiValue) bool {
	if m.Size() != other.Size() {
		return false
	}

	for i := range m.values {
		// TODO: add test to confirm this is correct!
		if m.values[i].EqualTo(other.values[i]) == False {
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
	if i >= len(m.values) {
		return NewUndefinedWithReasonf("out of bounds: trying to get arg #%d on MultiValue that has %d arguments", i, len(m.values))
	}

	return m.values[i]
}

func (m MultiValue) Size() int {
	return len(m.values)
}
