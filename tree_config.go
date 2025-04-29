package gal

import "strings"

// Variables holds the value of user-defined variables.
type Variables map[string]Value

func (v Variables) Get(name string) (Value, bool) {
	if v == nil {
		return nil, false
	}
	obj, ok := v[name]
	return obj, ok
}

// Functions holds the definition of user-defined functions.
type Functions map[string]FunctionalValue

func (f Functions) Get(name string) (FunctionalValue, bool) {
	if f == nil {
		return nil, false
	}
	obj, ok := f[name]
	return obj, ok
}

// Objects is a collection of Object's in the form of a map which keys are the name of the
// object and values are the actual Object's.
type Objects map[string]Object

// Get returns the Object of the specified name.
func (o Objects) Get(name string) (Object, bool) {
	if o == nil {
		return nil, false
	}
	obj, ok := o[name]
	return obj, ok
}

type treeConfig struct {
	variables Variables
	functions Functions
	objects   Objects
}

// Variable returns the value of the variable specified by name.
// TODO: add support for arrays and maps via `[...]`
// ...   NOTE: it may be more adequate to create a new `[]` operator.
// ...   This would also permit its use on any Value, including those returned from function calls.
// ...   We would likely need to create new types (unless MultiValue can work for this).
// ...   An awkward and visually less elegant option would be builtin functions such as GetIndex() (for arrays) and GetKey (for maps).
// ...................................................................
// ...................................................................
// ...   Perhaps this indicates that it's time to drop gal.Value   ...
// ...   and use native Go types and reflection?!?!                ...
// ...................................................................
// ...................................................................
func (tc treeConfig) Variable(name string) Value {
	if val, ok := tc.variables.Get(name); ok {
		return val
	}

	return NewUndefinedWithReasonf("error: unknown user-defined variable '%s'", name)
}

func (tc treeConfig) ObjectProperty(objProp ObjectProperty) Value {
	if obj, ok := tc.objects.Get(objProp.ObjectName); ok {
		return ObjectGetProperty(obj, objProp.PropertyName)
	}
	return NewUndefinedWithReasonf("error: object property '%s': unknown object", objProp.String())
}

// Function returns the function definition of the function of the specified name.
// This method is used to look up object methods and user-defined functions.
// Built-in functions are not looked up here, they are pre-populated at
// parsing time by the TreeBuilder.
func (tc treeConfig) Function(name string) FunctionalValue {
	splits := strings.Split(name, ".") // NOTE: strings.Split(name, ".") allocates a slice every call. strings.Count + strings.Index/LastIndex could avoid allocation in the common “no dot” path.
	if len(splits) == 2 {
		// look up the method in the user-provided objects
		tc.objectMethod(splits[0], splits[1])
		if obj, ok := tc.objects.Get(splits[0]); ok {
			// we ignore "ok" here because ObjectGetMethod will populate it with an Undefined.
			fv, _ := ObjectGetMethod(obj, splits[1])
			return fv
		}
		return func(...Value) Value {
			return NewUndefinedWithReasonf("error: object reference '%s' is not valid: unknown object or unknown method", name)
		}
	}

	if len(splits) >= 2 {
		// for expressions like `obj.a.b`, the tree should use a Variable or a Function to access `a` and
		//  then a Dot[Variable] / Dot[Function] to access `b`.
		return func(...Value) Value {
			return NewUndefinedWithReasonf("syntax error: object reference '%s' is not valid: too many dot accessors: max 1 permitted", name)
		}
	}

	// look up the function in the user-defined functions
	if val, ok := tc.functions.Get(name); ok {
		return val
	}

	return func(...Value) Value {
		return NewUndefinedWithReasonf("error: unknown user-defined function '%s'", name)
	}
}

// TODO: should this return a Function rather?
func (tc treeConfig) ObjectMethod(objMethod ObjectMethod) FunctionalValue {
	return tc.objectMethod(objMethod.ObjectName, objMethod.MethodName)
}

func (tc treeConfig) objectMethod(objectName, methodName string) FunctionalValue {
	if obj, ok := tc.objects.Get(objectName); ok {
		if fv, ok := ObjectGetMethod(obj, methodName); ok {
			return fv
		}
		return func(...Value) Value {
			return NewUndefinedWithReasonf("error: object '%s' method '%s': unknown or non-callable member (check if it has a pointer receiver)", objectName, methodName)
		}
	}

	return func(...Value) Value {
		return NewUndefinedWithReasonf("error: object '%s' method '%s': unknown object", objectName, methodName)
	}
}

type treeOption func(*treeConfig)

// WithVariables is a functional parameter for Tree evaluation.
// It provides user-defined variables.
func WithVariables(vars Variables) treeOption {
	return func(cfg *treeConfig) {
		cfg.variables = vars
	}
}

// WithFunctions is a functional parameter for Tree evaluation.
// It provides user-defined functions.
func WithFunctions(funcs Functions) treeOption {
	return func(cfg *treeConfig) {
		cfg.functions = funcs
	}
}

// WithObjects is a functional parameter for Tree evaluation.
// It provides user-defined Objects.
// These objects can carry both properties and methods that can be accessed
// by gal in place of variables and functions.
func WithObjects(objects Objects) treeOption {
	return func(cfg *treeConfig) {
		cfg.objects = objects
	}
}
