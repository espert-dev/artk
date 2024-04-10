// Package clone creates deep copies (as opposed to shallow copies).
package clone

import (
	"fmt"
	"reflect"
	"time"
)

// AsImmutable declares the type of the example value as immutable.
//
// Immutable types will not be deep copied.
// Only struct types can be assumed to be immutable.
// Cloning a value of an immutable type will never panic, even if it has
// unexported fields.
func AsImmutable(example any) {
	if example == nil {
		panic("example cannot be nil")
	}

	t := reflect.TypeOf(example)
	if t.Kind() != reflect.Struct {
		panic("only structs can be declared immutable")
	}

	key := typeIdOf(t)
	immutableStructTypes[key] = struct{}{}
}

// Of returns a clone (AKA deep copy) of the passed object.
//
// This function will panic if any of the below conditions are met:
//
//   - The value is or contains a channel.
//   - The value is or contains a function.
//   - The value is or contains an unsafe pointer.
//   - The value is or contains a struct with unexported fields, except if
//     that struct that has been marked with AsImmutable.
func Of[T any](x T) T {
	v := reflect.ValueOf(x)

	// Handle top-level nil interface values.
	if !v.IsValid() {
		var zero T
		return zero
	}

	k := cloner{
		knownPointers: make(map[uintptr]reflect.Value),
	}
	c := k.cloneAny(v)
	return c.Interface().(T)
}

type cloner struct {
	knownPointers map[uintptr]reflect.Value
}

func (k cloner) cloneAny(v reflect.Value) reflect.Value {
	switch kind := v.Kind(); kind {
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
		reflect.Complex128,
		reflect.String:
		// These kinds cannot introduce shallow copy.
		return v
	case reflect.Array:
		return k.cloneArray(v)
	case reflect.Interface:
		return k.cloneInterface(v)
	case reflect.Map:
		return k.cloneMap(v)
	case reflect.Pointer:
		return k.clonePointer(v)
	case reflect.Slice:
		return k.cloneSlice(v)
	case reflect.Struct:
		return k.cloneStruct(v)
	default:
		panic("unsupported kind: " + kind.String())
	}
}

func (k cloner) cloneArray(v reflect.Value) reflect.Value {
	// In arrays, the length is a property of the type, not the value.
	t := v.Type()
	l := t.Len()
	e := t.Elem()

	// Allocate memory for the cloned object.
	//
	// While it needs an addressable value, the particular address is
	// irrelevant. Hence, we immediately take Elem.
	c := reflect.New(reflect.ArrayOf(l, e)).Elem()

	// Recursively copy contents.
	// Avoid reflect.Copy, which would result in a shallow copy.
	for i := 0; i < l; i++ {
		x := v.Index(i)
		y := k.cloneAny(x)
		c.Index(i).Set(y)
	}

	return c
}

func (k cloner) cloneInterface(v reflect.Value) reflect.Value {
	// Nil values cannot lead to shallow copies.
	if v.IsNil() {
		return v
	}

	x := v.Elem()
	y := k.cloneAny(x)
	return y
}

func (k cloner) clonePointer(v reflect.Value) reflect.Value {
	// Nil values cannot lead to shallow copies.
	if v.IsNil() {
		return v
	}

	// Support cyclic data structures by reusing allocated pointers.
	p := v.Pointer()
	if v, ok := k.knownPointers[p]; ok {
		return v
	}

	// Allocate memory for the cloned object.
	x := v.Elem()
	t := x.Type()
	c := reflect.New(t)

	// Save pointer before recursing to prevent double allocation.
	k.knownPointers[p] = c

	// Recursively copy pointed-to element.
	y := k.cloneAny(x)
	c.Elem().Set(y)
	return c
}

func (k cloner) cloneMap(v reflect.Value) reflect.Value {
	// Nil values cannot lead to shallow copies.
	if v.IsNil() {
		return v
	}

	// Support cyclic data structures by reusing allocated pointers.
	p := v.Pointer()
	if v, ok := k.knownPointers[p]; ok {
		return v
	}

	// Allocate memory for the cloned object.
	t := v.Type()
	l := v.Len()
	c := reflect.MakeMapWithSize(t, l)

	// Save pointer before recursing to prevent double allocation.
	k.knownPointers[p] = c

	// Recursively copy contents.
	r := v.MapRange()
	for r.Next() {
		key := r.Key()
		x := r.Value()
		y := k.cloneAny(x)
		c.SetMapIndex(key, y)
	}

	return c
}

func (k cloner) cloneSlice(v reflect.Value) reflect.Value {
	// Nil values cannot lead to shallow copies.
	if v.IsNil() {
		return v
	}

	// Support cyclic data structures by reusing allocated pointers.
	p := v.Pointer()
	if v, ok := k.knownPointers[p]; ok {
		return v
	}

	// Allocate memory for the cloned object.
	//
	// Set capacity to length for memory efficiency. Note that in slices,
	// the latter is a property of the value, not the type.
	t := v.Type()
	l := v.Len()
	c := reflect.MakeSlice(t, l, l)

	// Save pointer before recursing to prevent double allocation.
	k.knownPointers[p] = c

	// Recursively copy contents.
	// Avoid reflect.Copy, which would result in a shallow copy.
	for i := 0; i < l; i++ {
		x := v.Index(i)
		y := k.cloneAny(x)
		c.Index(i).Set(y)
	}

	return c
}

func (k cloner) cloneStruct(v reflect.Value) reflect.Value {
	t := v.Type()

	// Handle immutable types such as time.Time, which contain unexported
	// fields but can be shared without concern.
	tid := typeIdOf(t)
	if _, ok := immutableStructTypes[tid]; ok {
		return v
	}

	// Allocate memory for the cloned object.
	//
	// We cannot use reflect.Zero because it returns a non-addressable
	// value, which would then fail when setting fields. We don't actually
	// care about the address, though, so we immediately take Elem.
	c := reflect.New(t).Elem()

	// NumField is a method in an interface, not a simple accessor, so we
	// cache it here since it's used twice.
	n := t.NumField()

	// The package reflect will panic on attempts to set an
	// unexported value. We preempt this situation to provide
	// a friendlier error message.
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if !f.IsExported() {
			panic(fmt.Sprintf(
				"struct has unexported fields: %v.%v",
				t.Name(),
				f.Name,
			))
		}
	}

	// Recursively clone fields.
	for i := 0; i < n; i++ {
		x := v.Field(i)
		y := k.cloneAny(x)
		c.Field(i).Set(y)
	}

	return c
}

// typeIdOf returns a global key for a type.
func typeIdOf(t reflect.Type) typeId {
	return typeId{
		PkgPath: t.PkgPath(),
		Name:    t.Name(),
	}
}

// init adds immutable types defined by the standard library.
func init() {
	AsImmutable(time.Time{})
}

// typeId represents the global key for a type.
type typeId struct {
	PkgPath string
	Name    string
}

// immutableStructTypes contains types that can be copied as values.
var immutableStructTypes = make(map[typeId]struct{})
