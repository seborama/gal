package gal

import (
	"fmt"
)

// TODO: could we use the same principle as Function.Receiver with Variable? Would it be elegant?
// ObjectProperty is a Tree entry that holds a reference of a user-defined object by name and the property to access on it.
// It is used to access a property on a user-defined object.
// It is a "cousin" of Variable, but for a property of a user-defined object.
type ObjectProperty struct {
	ObjectName   string
	PropertyName string
}

func NewObjectProperty(objectName, propertyName string) ObjectProperty {
	return ObjectProperty{
		ObjectName:   objectName,
		PropertyName: propertyName,
	}
}

//nolint:errcheck // life's too short to check for type assertion success here
func (o ObjectProperty) Calculate(val entry, op Operator, cfg *treeConfig) entry {
	rhsVal := cfg.ObjectProperty(o)
	if u, ok := rhsVal.(Undefined); ok {
		return u
	}

	if val == nil {
		return rhsVal
	}

	val = calculate(val.(Value), op, rhsVal)

	return val
}

func (o ObjectProperty) String() string {
	return fmt.Sprintf("%s.%s", o.ObjectName, o.PropertyName)
}
