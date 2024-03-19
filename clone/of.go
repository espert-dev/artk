// Package clone creates deep copies (as opposed to shallow copies).
package clone

import (
	"reflect"
)

// Of returns a clone of the passed object. Private fields are not cloned.
func Of[T any](x T) T {
	v := reflect.ValueOf(x)
	c := cloneAny(v)
	return c.Interface().(T)
}

func cloneAny(v reflect.Value) reflect.Value {
	// Must cover all values of reflect.Kind
	switch v.Type().Kind() {
	case reflect.Invalid:
		panic("invalid kind")
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		// Returning the value is enough for these value types.
		return v
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Pointer:
	case reflect.Slice:
		return cloneSlice(v)
	case reflect.String:
	case reflect.Struct:
	case reflect.UnsafePointer:
	}

	panic("unsupported kind")
}

func cloneSlice(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	// Set capacity to length.
	l := v.Len()
	c := reflect.MakeSlice(v.Type(), l, l)

	// We avoid reflect.Copy, which would result in a shadow copy.
	for i := 0; i < l; i++ {
		e := v.Index(i)
		c.Index(i).Set(e)
	}

	return c
}
