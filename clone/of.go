// Package clone creates deep copies (as opposed to shallow copies).
package clone

import (
	"fmt"
	"reflect"
)

// Of returns a clone of the passed object. Private fields are not cloned.
func Of[T any](x T) T {
	v := reflect.ValueOf(x)

	// Can handle nil interface values.
	if !v.IsValid() {
		var zero T
		return zero
	}

	c := cloneAny(v)
	return c.Interface().(T)
}

func cloneAny(v reflect.Value) reflect.Value {
	// Must cover all values of reflect.Kind.
	switch v.Kind() {
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
		return cloneArray(v)
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
		panic("it's an interface")
	case reflect.Map:
		return cloneMap(v)
	case reflect.Pointer:
		return clonePointer(v)
	case reflect.Slice:
		return cloneSlice(v)
	case reflect.String:
		// Strings are immutable in Go.
		return v
	case reflect.Struct:
		return cloneStruct(v)
	case reflect.UnsafePointer:
	}

	panic("unsupported kind")
}

func cloneArray(v reflect.Value) reflect.Value {
	// In arrays, the length is a property of the type, not the value.
	l := v.Type().Len()
	t := v.Type().Elem()
	c := reflect.New(reflect.ArrayOf(l, t))

	// We avoid reflect.Copy, which would result in a shadow copy.
	for i := 0; i < l; i++ {
		x := v.Index(i)
		y := cloneAny(x)
		c.Elem().Index(i).Set(y)
	}

	return c.Elem()
}

func clonePointer(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	x := v.Elem()
	y := cloneAny(x)

	t := x.Type()
	c := reflect.New(t)
	c.Elem().Set(y)
	return c
}

func cloneMap(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	t := v.Type()
	l := v.Len()
	c := reflect.MakeMapWithSize(t, l)

	r := v.MapRange()
	for r.Next() {
		k := r.Key()
		x := r.Value()
		y := cloneAny(x)
		c.SetMapIndex(k, y)
	}

	return c
}

func cloneSlice(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	// Set capacity to length.
	// Note that in slices, this is a property of the value, not the type.
	t := v.Type()
	l := v.Len()
	c := reflect.MakeSlice(t, l, l)

	// We avoid reflect.Copy, which would result in a shadow copy.
	for i := 0; i < l; i++ {
		x := v.Index(i)
		y := cloneAny(x)
		c.Index(i).Set(y)
	}

	return c
}

func cloneStruct(v reflect.Value) reflect.Value {
	t := v.Type()

	// We cannot use reflect.Zero because it returns a non-addressable
	// value, which would then fail when setting fields.
	c := reflect.New(t)

	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)

		// The package reflect will panic on attempts to set an
		// unexported value. We preempt this situation to provide
		// a friendlier error message.
		if !f.IsExported() {
			panic(fmt.Sprintf(
				"struct has unexported fields: %v.%v",
				t.Name(),
				f.Name,
			))
		}

		x := v.Field(i)
		y := cloneAny(x)
		c.Elem().Field(i).Set(y)
	}

	return c.Elem()
}
