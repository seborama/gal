package gal

import (
	"fmt"
	"reflect"
)

// TODO: implement support for nested structs?
func ObjectGetProperty(obj Object, name string) (Value, bool) {
	if obj == nil {
		return NewUndefinedWithReasonf("object is nil"), false
	}

	// Use the reflect.ValueOf function to get the value of the struct
	v := reflect.ValueOf(obj)
	if v.IsZero() || !v.IsValid() {
		return NewUndefinedWithReasonf("object is nil, not a Go value or invalid"), false
	}

	// Use reflect.TypeOf to get the type of the struct
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.IsZero() || !v.IsValid() {
			return NewUndefinedWithReasonf("object interface is nil, not a Go value or invalid"), false
		}

		t = t.Elem()
		if v.IsZero() || !v.IsValid() {
			return NewUndefinedWithReasonf("object interface is nil, not a Go value or invalid"), false
		}
	}

	if t.Kind() != reflect.Struct {
		// TODO: we only support `struct` for now. Perhaps simple types (int, float, etc) are a worthwhile enhancement?
		return NewUndefinedWithReasonf("object is '%s' but only 'struct' and '*struct' are currently supported", t.Kind()), false
	}

	typeName := t.Name()

	// Iterate over the fields of the struct
	for i := 0; i < v.NumField(); i++ {
		vName := v.Type().Field(i).Name
		vType := v.Type().Field(i).Type.Name()
		vValueI := v.Field(i).Interface()
		if vName == name {
			if vValue, ok := vValueI.(Value); ok {
				return vValue, true
			} else {
				switch vValueIType := vValueI.(type) {
				case int:
					return NewNumberFromInt(int64(vValueIType)), true
				case int32:
					return NewNumberFromInt(int64(vValueIType)), true
				case int64:
					return NewNumberFromInt(vValueIType), true
				case uint:
					return NewNumberFromInt(int64(vValueIType)), true
				case uint32:
					return NewNumberFromInt(int64(vValueIType)), true
				case uint64:
					n, err := NewNumberFromString(fmt.Sprintf("%d", vValueIType))
					if err != nil {
						return NewUndefinedWithReasonf("value '%d' of property '%s:%s' cannot be converted to a Number", vValueIType, typeName, name), false
					}
					return n, true
				case float32: // this will commonly suffer from floating point issues
					return NewNumberFromFloat(float64(vValueIType)), true
				case float64:
					return NewNumberFromFloat(vValueIType), true
				case string:
					return NewString(vValueIType), true
				case bool:
					return NewBool(vValueIType), true
				default:
					return NewUndefinedWithReasonf("property '%s:%s' is of type '%s', not a gal.Value", typeName, name, vType), false
				}
			}
		}
	}

	return NewUndefinedWithReasonf("property '%s:%s' does not exist on object", typeName, name), false
}
