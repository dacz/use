package use

import "reflect"

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Interface:
		return v.IsZero() || v.IsNil()
	default:
		return false
	}
}

func derefValue(vv reflect.Value) (indirect bool, v reflect.Value) {
	if vv.Kind() == reflect.Ptr {
		return true, vv.Elem()
	}
	return false, vv
}

func derefType(vt reflect.Type) (indirect bool, t reflect.Type) {
	if vt.Kind() == reflect.Ptr {
		return true, vt.Elem()
	}
	return false, vt
}

func typesMatch(outFt, inpFt reflect.Type) bool {
	if outFt == inpFt {
		return true
	}

	if outFt.Kind() == reflect.Ptr && outFt.Elem() == inpFt {
		return true
	}

	if inpFt.Kind() == reflect.Ptr && inpFt.Elem() == outFt {
		return true
	}

	return false
}

// containsStructOrPtrToStruct does not suport double (or more) references (e.g. **SomeStruct). It is often a mistake to use them.
func containsStructOrPtrToStruct(t reflect.Type) bool {
	if t.Kind() == reflect.Struct {
		return true
	}

	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return true
	}

	return false
}

func asRef[T any](s T) *T {
	return &s
}

func addToFields(parentFieldName, fieldName string) string {
	if parentFieldName == "" {
		return fieldName
	}
	return parentFieldName + "." + fieldName
}
