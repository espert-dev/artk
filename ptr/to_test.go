package ptr_test

import (
	"artk.dev/ptr"
	"reflect"
	"testing"
	"unsafe"
)

func TestTo_comparable(t *testing.T) {
	// Avoid a traditional table test so that `T` isn't always `any`.
	// We need a test case for each value of reflect.Kind except for
	// Invalid.
	//
	// For extra safety, we avoid checking against the zero values of the
	// types, since those could hide a lack of value initialisation.

	// Booleans.
	testComparable(t, true)

	// Numbers.
	const x = 42
	testComparable(t, x)
	testComparable(t, int8(x))
	testComparable(t, int16(x))
	testComparable(t, int32(x))
	testComparable(t, int64(x))
	testComparable(t, uint(x))
	testComparable(t, uint8(x))
	testComparable(t, uint16(x))
	testComparable(t, uint32(x))
	testComparable(t, uint64(x))
	testComparable(t, uintptr(x))
	testComparable(t, float32(x))
	testComparable(t, float64(x))
	testComparable(t, complex64(x))
	testComparable(t, complex128(x))

	// Comparable data structures.
	p := new(int)
	*p = x

	testComparable(t, [1]int{x})
	testComparable(t, make(chan int))
	testComparable(t, "foo")
	testComparable(t, p)
	testComparable(t, struct{ Value int }{Value: x})
	testComparable(t, unsafe.Pointer(&struct{}{}))
}

func TestTo_Interface(t *testing.T) {
	// reflect.TypeOf returns the concrete type of the value, not the
	// interface used to describe it. This means that while `any` is
	// comparable, we still need a special case for it.
	testComparable(t, any(true))
}

func TestTo_Map(t *testing.T) {
	m := map[string]int{"foo": 42}
	p := ptr.To(m)
	requireValidPointer(t, m, p)
	if !reflect.DeepEqual(m, *p) {
		t.Errorf("expected %v, got %v", m, *p)
	}
}

func TestTo_Slice(t *testing.T) {
	s := []int{42}
	p := ptr.To(s)
	requireValidPointer(t, s, p)
	if !reflect.DeepEqual(s, *p) {
		t.Errorf("expected %v, got %v", s, *p)
	}
}

func TestTo_Func(t *testing.T) {
	// Functions are non-comparable, so they need their own special case.
	// We can only observe some sense of function equality via
	// their side effects.
	sideEffect := false
	fn := func() {
		sideEffect = true
	}

	ptrToFn := ptr.To(fn)
	requireValidPointer(t, fn, ptrToFn)

	// Observe the side effect.
	if sideEffect {
		t.Fatal("sideEffect was unexpectedly true")
	}
	(*ptrToFn)()
	if !sideEffect {
		t.Fatal("sideEffect was unexpectedly false")
	}
}

func testComparable[T comparable](t *testing.T, value T) {
	t.Helper()
	t.Run(reflect.TypeOf(value).String(), func(t *testing.T) {
		pointer := ptr.To(value)
		requireValidPointer(t, value, pointer)
		if got := *pointer; got != value {
			t.Errorf("expected value %v, got %v", value, got)
		}
	})
}

func requireValidPointer[T any](t *testing.T, _ T, pointer *T) {
	// The correct type of the pointer is enforced by Go's type system.
	// We only pass the value T to get a compilation error if the generic
	// constraints cannot be met. As a bonus, this supports `any` and
	// comparing against nil.
	if pointer == nil {
		t.Fatal("unexpectedly returned a nil pointer")
	}
}
