package gal

import (
	"fmt"
)

// ObjectMethod is a Tree entry that holds a reference of a user-defined object by name and the method to call on it.
// It is used to call a method on a user-defined object.
// It is a "cousin" of Function, but for a method of a user-defined object.
type ObjectMethod struct {
	ObjectName string
	MethodName string
	Args       []Tree
}

func NewObjectMethod(objectName, propertyName string, args ...Tree) ObjectMethod {
	return ObjectMethod{
		ObjectName: objectName,
		MethodName: propertyName,
		Args:       args,
	}
}

//nolint:errcheck // life's too short to check for type assertion success here
func (om ObjectMethod) Calculate(val entry, op Operator, cfg *treeConfig) entry {
	// attempt to get body of a user-provided object's method.
	bodyFn := cfg.ObjectMethod(om)

	fn := NewFunction(om.MethodName, bodyFn, om.Args...)

	rhsVal := fn.Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables), WithObjects(cfg.objects))
	if u, ok := rhsVal.(Undefined); ok {
		return u
	}

	if val == nil {
		return rhsVal
	}

	val = calculate(val.(Value), op, rhsVal)

	return val
}

func (om ObjectMethod) String() string {
	return fmt.Sprintf("%s.%s", om.ObjectName, om.MethodName)
}
