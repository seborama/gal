package gal

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// Object holds user-defined objects that can carry properties and functions that may be
// referenced within a gal expression during evaluation.
type Object any

// TODO: implement support for nested structs?
func ObjectGetProperty(obj Object, name string) (Value, bool) { //nolint: gocognit, gocyclo, cyclop
	if obj == nil {
		return NewUndefinedWithReasonf("object is nil"), false
	}

	// Use the reflect.ValueOf function to get the value of the struct
	v := reflect.ValueOf(obj)
	if !v.IsValid() {
		return NewUndefinedWithReasonf("object is nil, not a Go value or invalid"), false
	}

	// Use reflect.TypeOf to get the type of the struct
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if !v.IsValid() {
			return NewUndefinedWithReasonf("object interface is nil, not a Go value or invalid"), false
		}

		t = t.Elem()
		if !v.IsValid() {
			return NewUndefinedWithReasonf("object interface is nil, not a Go value or invalid"), false
		}
	}

	// TODO: we only support `struct` for now. Perhaps simple types (int, float, etc) are a worthwhile enhancement?
	if t.Kind() != reflect.Struct {
		return NewUndefinedWithReasonf("object is '%s' but only 'struct' and '*struct' are currently supported", t.Kind()), false
	}

	vValue := v.FieldByName(name)
	if !vValue.IsValid() {
		return NewUndefinedWithReasonf("property '%T:%s' does not exist on object", obj, name), false
	}

	slog.Debug("ObjectGetProperty", "vValue.Kind", vValue.Kind().String(), "name", name, "vValue", vValue)
	if vValue.Kind() == reflect.Struct {
		// TODO: do not hard-code "Age", enable a means to continue parsing the name.
		// ...   NOTE: it may be needed to create a wrapper over ObjectGetProperty and ObjectGetMethod that
		// ...   can continue iterating thought the '.' separated parts of the name and depending on the presence
		// ...   of () at the end of the part, call ObjectGetProperty or ObjectGetMethod accordingly
		// ...   NOTE: Instead of using functionEntryKind / variableEntryKind, we might want use objectEntryKind.
		// ...   This would allow a more universal parsing of expressions like:
		// ...   aCar.Tyres[2],Vertices(4).Length.InInches()
		// ...   This should also remove some burden away from gal's core TreeBuild and Tree and provide
		// ...   decoupling / separation of concerns between them and the evalution of "object" expressions.
		return ObjectGetProperty(vValue.Interface(), "Age")
	}

	galValue, err := goAnyToGalType(vValue.Interface())
	if err != nil {
		return NewUndefinedWithReasonf("object::%T:%s - %s", obj, name, err.Error()), false
	}
	return galValue, true
}

func ObjectGetMethod(obj Object, name string) (FunctionalValue, bool) { //nolint: cyclop
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

	methodValue := value.MethodByName(name)
	if !methodValue.IsValid() {
		return func(...Value) Value {
			return NewUndefinedWithReasonf("type '%T' does not have a method '%s' (check if it has a pointer receiver)", obj, name)
		}, false
	}

	methodType := methodValue.Type()
	numParams := methodType.NumIn()

	var fn FunctionalValue = func(args ...Value) (retValue Value) {
		if len(args) != numParams {
			return NewUndefinedWithReasonf("invalid function call - object::%T:%s - wants %d args, received %d instead", obj, name, numParams, len(args))
		}

		callArgs := lo.Map(args, func(item Value, index int) reflect.Value {
			paramType := methodType.In(index)

			switch paramType.Kind() {
			// TODO: continue with more "case"'s
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

		out := methodValue.Call(callArgs)
		if len(out) != 1 {
			return NewUndefinedWithReasonf("invalid function call - object::%T:%s - must return 1 value, returned %d instead", obj, name, len(out))
		}

		retValue, err := goAnyToGalType(out[0].Interface())
		if err != nil {
			return NewUndefinedWithReasonf("object::%T:%s - %s", obj, name, err.Error())
		}
		return
	}

	return fn, true
}

// attempt to convert a Go 'any' type to an equivalent gal.Value
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
		t := reflect.TypeOf(value)
		slog.Debug("goAnyToGalType", "t.Kind", t.Kind().String())
		return nil, errors.Errorf("type '%T' cannot be mapped to gal.Value", typedValue)
	}
}
