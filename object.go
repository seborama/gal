package gal

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type DotFunction struct{ Function }

func (df DotFunction) Calculate(val entry, cfg *treeConfig) entry {
	if df.BodyFn != nil {
		// NOTE: this could be supported but it would turn the object into a prototype model e.g. like JavaScript
		return NewUndefinedWithReasonf("internal error: DotFunction for '%s': BodyFn is not empty: this indicates the object's method was confused for a build-in function", df.Name)
	}

	var receiver any

	// as this is an object function accessor, we need to get the object first: it is the LHS currently held in val
	receiver, ok := val.(Value)
	if !ok {
		return NewUndefinedWithReasonf("syntax error: DotFunction called on non-object: [object: '%T'] [member: '%s'] (check if the receiver is nil)", val, df.Name)
	}

	// if the object is a ObjectValue, we need to get the underlying object
	// ObjectValue is a wrapper for "general" objects (i.e. non-gal.Value objects)
	// By Object, we mean a Go struct, a pointer to a struct or a Go interface.
	objVal, ok := receiver.(ObjectValue)
	if ok {
		receiver = objVal.Object
	}

	// now, we can get the method from the object
	vFv, ok := ObjectGetMethod(receiver, df.Name)
	if ok {
		df.BodyFn = vFv
		rhsVal := df.Eval(WithFunctions(cfg.functions), WithVariables(cfg.variables), WithObjects(cfg.objects))
		if u, ok := rhsVal.(Undefined); ok {
			return u
		}

		return rhsVal
	}

	return vFv // this will already be an Undefined type.
}

type DotVariable struct{ Variable }

func (dv DotVariable) Calculate(val entry) entry {
	var receiver any

	// as this is an object property accessor, we need to get the object first: it is the LHS currently held in val
	receiver, ok := val.(Value)
	if !ok {
		return NewUndefinedWithReasonf("syntax error: object accessor [Variable] called on non-object: [object: '%T'] [member: '%s'] (check if the receiver is nil)", fmt.Sprintf("%T", val), dv.Name)
	}

	// if the object is a ObjectValue, we need to get the underlying object
	// ObjectValue is a wrapper for "general" objects (i.e. non-gal.Value objects)
	// By Object, we mean a Go struct, a pointer to a struct or a Go interface.
	objVal, ok := receiver.(ObjectValue)
	if ok {
		receiver = objVal.Object
	}

	// now, we can get the property from the object
	rhsVal := ObjectGetProperty(receiver, dv.Name)

	return rhsVal
}

// Object holds objects that carry properties and methods:
//   - user-defined objects that may be referenced within a gal expression during evaluation.
//   - general purpose Go types that have properties and methods.
//     These are provided by user-defined objects via their properties and methods return values.
type Object any

// ObjectValue is a "bridge" beween a non-Value object and Value.
// This is useful for object accessors that return a non-value.
// While we cannot perform any Value operations on such objects,
// ObjectValue allows to keep "traversing" the objects with the Dot
// operator until (hopefully) we end with a Value.
type ObjectValue struct {
	Object any
	Undefined
}

func (o ObjectValue) String() string {
	return fmt.Sprintf("ObjectValue(%T)", o.Object)
}

// ObjectGetProperty returns the value of the property with the given name from the object.
func ObjectGetProperty(obj Object, name string) Value {
	if obj == nil {
		return NewUndefinedWithReasonf("object is nil for type '%T'", obj)
	}

	// Use the reflect.ValueOf function to get the value of the struct
	v := reflect.ValueOf(obj)
	if !v.IsValid() {
		return NewUndefinedWithReasonf("object is nil, not a Go value or invalid")
	}

	// Use reflect.TypeOf to get the type of the struct
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if !v.IsValid() {
			return NewUndefinedWithReasonf("object interface is nil, not a Go value or invalid")
		}

		t = t.Elem()
		if !v.IsValid() {
			return NewUndefinedWithReasonf("object interface is nil, not a Go value or invalid")
		}
	}

	// TODO: we only support `struct` receivers for now. Perhaps simple types (int, float, etc) are worthwhile an enhancement?
	if t.Kind() != reflect.Struct {
		return NewUndefinedWithReasonf("object is '%s' but only 'struct' and '*struct' are currently supported", t.Kind())
	}

	fieldReflectValue := v.FieldByName(name)
	if !fieldReflectValue.IsValid() {
		return NewUndefinedWithReasonf("property '%T:%s' does not exist on object", obj, name)
	}

	galValue, err := goAnyToGalType(fieldReflectValue.Interface())
	if err != nil {
		// allow support for other types to be accessed by Method or Property via
		//  an object accessor (i.e. DotVariable or DotFunction).
		t := fieldReflectValue.Type()
		switch t.Kind() {
		case reflect.Interface:
			if t.NumMethod() > 0 {
				// allow support for (non-empty) interfaces
				return ObjectValue{Object: fieldReflectValue.Interface()}
			}
		case reflect.Struct: // TODO: (!!) incomplete code: see ObjectGetProperty to handle `*struct` scenario.
			// allow support for struct types
			return ObjectValue{Object: fieldReflectValue.Interface()}
		}

		return NewUndefinedWithReasonf("object::%T:%s - %s", obj, name, err.Error())
	}

	return galValue
}

// ObjectGetMethod returns a closure that can be called with the method's arguments.
func ObjectGetMethod(obj Object, name string) (FunctionalValue, bool) {
	if obj == nil {
		return func(...Value) Value {
			return NewUndefinedWithReasonf("object is nil for type '%T'", obj)
		}, false
	}

	value := reflect.ValueOf(obj)
	if !value.IsValid() {
		return func(...Value) Value {
			return NewUndefinedWithReasonf("object type '%T' is not valid", obj)
		}, false
	}

	methodReflectValue := value.MethodByName(name)
	if !methodReflectValue.IsValid() {
		return func(...Value) Value {
			return NewUndefinedWithReasonf("error: object type '%T' does not have a method '%s' (check if it has a pointer receiver)", obj, name)
		}, false
	}

	methodType := methodReflectValue.Type()
	numParams := methodType.NumIn()

	var closureFn FunctionalValue = func(args ...Value) (retValue Value) {
		if len(args) != numParams {
			return NewUndefinedWithReasonf("invalid function call - object::%T:%s - wants %d args, received %d instead", obj, name, numParams, len(args))
		}

		//nolint:gosec // ignoring mem overflow conversion
		// for functions that requires non-gal.Value parameters, attempt to map such gal.Value params to
		// what the function param dictates.
		// E.g.: if an object method has a signature of func (int), we will attempt to map the gal.Value to
		// a Numberer, then extract an Int64 from the Number() and finally map it to an "int".
		callArgs := lo.Map(args, func(item Value, index int) reflect.Value {
			paramType := methodType.In(index)

			//nolint:errcheck // life's too short to check for type assertion success here
			switch paramType.Kind() {
			case reflect.Int:
				return reflect.ValueOf(int(item.(Numberer).Number().Int64()))
			case reflect.Int32:
				return reflect.ValueOf(int32(item.(Numberer).Number().Int64()))
			case reflect.Int64:
				return reflect.ValueOf(item.(Numberer).Number().Int64())
			case reflect.Uint:
				return reflect.ValueOf(uint(item.(Numberer).Number().Int64()))
			case reflect.Uint32:
				return reflect.ValueOf(uint32(item.(Numberer).Number().Int64()))
			case reflect.Uint64:
				n, err := strconv.ParseUint(item.(Stringer).AsString().RawString(), 10, 64)
				if err != nil {
					panic(err) // no other safe way
				}
				return reflect.ValueOf(n)
			case reflect.Float32:
				return reflect.ValueOf(float32(item.(Numberer).Number().Float64()))
			case reflect.Float64:
				return reflect.ValueOf(item.(Numberer).Number().Float64())
			case reflect.String:
				return reflect.ValueOf(item.(Stringer).AsString().RawString())
			case reflect.Bool:
				return reflect.ValueOf(item.(Booler).Bool().value)
			default:
				return reflect.ValueOf(item)
			}
		})

		defer func() {
			if r := recover(); r != nil {
				retValue = NewUndefinedWithReasonf("invalid function call - object::%T:%s - invalid argument type passed to function - %v", obj, name, r)
				return
			}
		}()

		out := methodReflectValue.Call(callArgs)
		if len(out) != 1 {
			return NewUndefinedWithReasonf("invalid function call - object::%T:%s - must return 1 value, returned %d instead", obj, name, len(out))
		}

		retValue, err := goAnyToGalType(out[0].Interface())
		if err != nil {
			// allow support for other types to be accessed by Method or Property via
			//  an objectAccessorEntryKind (i.e. DotVariable or DotFunction).
			t := out[0].Type()
			switch t.Kind() {
			case reflect.Interface:
				if t.NumMethod() > 0 {
					// allow support for (non-empty) interfaces
					return ObjectValue{Object: out[0].Interface()}
				}
			case reflect.Struct: // TODO: (!!) incomplete code: see ObjectGetProperty to handle `*struct` scenario.
				// allow support for struct types
				return ObjectValue{Object: out[0].Interface()}
			}

			return NewUndefinedWithReasonf("object::%T:%s - %s", obj, name, err.Error())
		}
		return retValue
	}

	return closureFn, true
}

// attempt to convert a Go 'any' type to an equivalent gal.Value
//
//nolint:gosec // ignoring overflow conversion
func goAnyToGalType(value any) (Value, error) {
	switch typedValue := value.(type) {
	case Value:
		return typedValue, nil
	case int:
		return NewNumberFromInt(int64(typedValue)), nil
	case int32:
		return NewNumberFromInt(int64(typedValue)), nil
	case int64:
		return NewNumberFromInt(typedValue), nil
	case uint:
		return NewNumberFromInt(int64(typedValue)), nil
	case uint32:
		return NewNumberFromInt(int64(typedValue)), nil
	case uint64:
		n, err := NewNumberFromString(fmt.Sprintf("%d", typedValue))
		if err != nil {
			return nil, errors.Errorf("value uint64(%d) cannot be converted to a Number", typedValue)
		}
		return n, nil
	case float32: // this will commonly suffer from floating point issues
		return NewNumberFromFloat(float64(typedValue)), nil
	case float64:
		return NewNumberFromFloat(typedValue), nil
	case string:
		return NewString(typedValue), nil
	case bool:
		return NewBool(typedValue), nil
	default:
		return nil, errors.Errorf("type '%T' cannot be mapped to gal.Value", typedValue)
	}
}
