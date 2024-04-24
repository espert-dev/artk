package clone_test

import (
	"artk.dev/clone"
	"artk.dev/internal/assert"
	"artk.dev/internal/require"
	"artk.dev/typetraits"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

// Derived types to ensure that cloning works for any type in the kind.
type (
	boolType       bool
	intType        int
	int8Type       int8
	int16Type      int16
	int32Type      int32
	int64Type      int64
	uintType       uint
	uint8Type      uint8
	uint16Type     uint16
	uint32Type     uint32
	uint64Type     uint64
	uintptrType    uintptr
	float32Type    float32
	float64Type    float64
	complex64Type  complex64
	complex128Type complex128
	stringType     string
)

func TestAsImmutableType_example_cannot_be_nil(t *testing.T) {
	defer assert.PanicBecause(t, "example cannot be nil")
	clone.AsImmutableType(nil)
}

func TestAsImmutableType_example_must_be_of_a_struct_type(t *testing.T) {
	for _, v := range []any{
		false,
		boolType(false),
		0,
		intType(0),
		int8(0),
		int8Type(0),
		int16(0),
		int16Type(0),
		int32(0),
		int32Type(0),
		int64(0),
		int64Type(0),
		uint(0),
		uintType(0),
		uint8(0),
		uint8Type(0),
		uint16(0),
		uint16Type(0),
		uint32(0),
		uint32Type(0),
		uint64(0),
		uint64Type(0),
		uintptr(0),
		uintptrType(0),
		float32(0),
		float32Type(0),
		float64(0),
		float64Type(0),
		complex64(0),
		complex64Type(0),
		complex128(0),
		complex128Type(0),
		[0]bool{},
		func() {},
		chan struct{}(nil),
		&time.Time{},
		map[bool]struct{}(nil),
		[]bool(nil),
		"",
		stringType(""),
		unsafe.Pointer(nil),
	} {
		t.Run(reflect.TypeOf(v).String(), func(t *testing.T) {
			const why = "only structs can be declared immutable"
			defer assert.PanicBecause(t, why)
			clone.AsImmutableType(v)
		})
	}
}

func TestAsImmutableType_is_idempotent(t *testing.T) {
	defer assert.NoPanic(t)
	clone.AsImmutableType(time.Time{})
	clone.AsImmutableType(time.Time{})
}

func TestOf_supports_booleans(t *testing.T) {
	testValues(t, []bool{false, true})
	testValues(t, []boolType{false, true})
}

func TestOf_supports_numbers(t *testing.T) {
	// Don't test with NaN because comparison is always false.
	pInf := math.Inf(1)
	nInf := math.Inf(-1)

	testValues(t, []int{-1, 0, 1})
	testValues(t, []intType{-1, 0, 1})
	testValues(t, []int8{-1, 0, 1})
	testValues(t, []int8Type{-1, 0, 1})
	testValues(t, []int16{-1, 0, 1})
	testValues(t, []int16Type{-1, 0, 1})
	testValues(t, []int32{-1, 0, 1})
	testValues(t, []int32Type{-1, 0, 1})
	testValues(t, []int64{-1, 0, 1})
	testValues(t, []int64Type{-1, 0, 1})
	testValues(t, []uint{0, 1})
	testValues(t, []uintType{0, 1})
	testValues(t, []uint8{0, 1})
	testValues(t, []uint8Type{0, 1})
	testValues(t, []uint16{0, 1})
	testValues(t, []uint16Type{0, 1})
	testValues(t, []uint32{0, 1})
	testValues(t, []uint32Type{0, 1})
	testValues(t, []uint64{0, 1})
	testValues(t, []uint64Type{0, 1})
	testValues(t, []uintptr{0, 1})
	testValues(t, []uintptrType{0, 1})
	testValues(t, []float32{0, 1})
	testValues(t, []float32Type{0, 1, float32Type(pInf), float32Type(nInf)})
	testValues(t, []float64{0, 1})
	testValues(t, []float64Type{0, 1, float64Type(pInf), float64Type(nInf)})
	testValues(t, []complex64{0, 1, complex(0, 1), complex(1, 1)})
	testValues(t, []complex64Type{0, 1, complex(0, 1), complex(1, 1)})
	testValues(t, []complex128{0, 1, complex(0, 1), complex(1, 1)})
	testValues(t, []complex128Type{0, 1, complex(0, 1), complex(1, 1)})
}

func testValues[T comparable](t *testing.T, values []T) {
	t.Helper()

	t.Run(reflect.TypeOf(values[0]).String(), func(t *testing.T) {
		for _, v := range values {
			t.Run(fmt.Sprintf("%v", v), func(t *testing.T) {
				got := clone.Of(v)
				assert.Equal(t, v, got)
			})
		}
	})
}

func TestOf_supports_arrays(t *testing.T) {
	testArray(t, [2]bool{})
	testArray(t, [2]bool{})
	testArray(t, [2]bool{false})
	testArray(t, [2]bool{true})
	testArray(t, [2]bool{false, true})
	testArray(t, [2]boolType{})
	testArray(t, [2]boolType{})
	testArray(t, [2]boolType{false})
	testArray(t, [2]boolType{true})
	testArray(t, [2]boolType{false, true})
}

// We have no way of parametrizing the array length with generics,
// so we picked an arbitrary but useful value.
func testArray[T comparable](t *testing.T, array [2]T) {
	t.Helper()
	t.Run(fmt.Sprintf("%T%v", array, array), func(t *testing.T) {
		c := clone.Of(array)
		assert.DeepEqual(t, array, c)
	})
}

func TestOf_supports_interfaces(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var v any
		c := clone.Of(v)
		assert.Equal(t, v, c)
	})
	t.Run("nested nil interface", func(t *testing.T) {
		type NestedInterface struct {
			Any any
		}
		v := NestedInterface{Any: nil}

		c := clone.Of(v)
		assert.Equal(t, v, c)
	})
	t.Run("not nil", func(t *testing.T) {
		var v any = true
		c := clone.Of(v)
		assert.Equal(t, v, c)
	})
	t.Run("nil interface", func(t *testing.T) {
		var v any = (*bool)(nil)
		c := clone.Of(v)
		assert.Equal(t, v, c)
	})
	t.Run("nested nil interface", func(t *testing.T) {
		type NestedInterface struct {
			Any any
		}
		v := NestedInterface{Any: (*bool)(nil)}

		c := clone.Of(v)
		assert.Equal(t, v, c)
	})
}

func TestOf_supports_pointers(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var v *int

		c := clone.Of(v)
		assert.Nil(t, c)
	})
	t.Run("not nil", func(t *testing.T) {
		v := new(int)
		*v = 1

		c := clone.Of(v)
		assert.NotSame(t, c, v)
		assert.Equal(t, *c, *v)
	})
}

func TestOf_supports_acyclic_maps(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input map[int]bool
	}{
		{
			name:  "nil",
			input: nil,
		},
		{
			name:  "empty",
			input: map[int]bool{},
		},
		{
			name:  "one",
			input: map[int]bool{0: true},
		},
		{
			name:  "two",
			input: map[int]bool{0: true, 1: false},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := clone.Of(tt.input)
			assert.NotSame(t, tt.input, c)
			assert.DeepEqual(t, tt.input, c)
		})
	}
}

func TestOf_supports_cyclic_maps(t *testing.T) {
	v := make(map[string]any)
	v["next"] = v

	c := clone.Of(v)
	assert.NotSame(t, v, c)
	assert.Same(t, v["next"], any(v))
}

func TestOf_supports_acyclic_slices(t *testing.T) {
	testSlice(t, []bool(nil))
	testSlice(t, []bool{})
	testSlice(t, []bool{false})
	testSlice(t, []bool{true})
	testSlice(t, []bool{false, true})
	testSlice(t, []boolType(nil))
	testSlice(t, []boolType{})
	testSlice(t, []boolType{false})
	testSlice(t, []boolType{true})
	testSlice(t, []boolType{false, true})
}

func testSlice[T comparable](t *testing.T, slice []T) {
	t.Helper()

	var name string
	if slice == nil {
		name = fmt.Sprintf("%T(nil)", slice)
	} else {
		name = fmt.Sprintf("%T%v", slice, slice)
	}

	t.Run(name, func(t *testing.T) {
		c := clone.Of(slice)
		assert.DeepEqual(t, slice, c)

		// Sharing is never an issue if the slices are empty.
		if len(slice) != 0 {
			assert.NotSame(t, slice, c)
		}
	})
}

func TestOf_supports_cyclic_slices(t *testing.T) {
	v := make([]any, 1)
	v[0] = v

	c := clone.Of(v)
	assert.NotSame(t, v, c)
	assert.Same(t, v[0], any(v))
}

func TestOf_supports_strings(t *testing.T) {
	testDeepEqual(t, "")
	testDeepEqual(t, stringType(""))
	testDeepEqual(t, "foo")
	testDeepEqual(t, stringType("foo"))
}

func TestOf_supports_acyclic_structs_without_unexported_fields(t *testing.T) {
	testDeepEqual(t, struct{}{})

	type IntStruct struct {
		X int
	}
	testDeepEqual(t, IntStruct{X: 1})

	type Point struct {
		X, Y, Z float64
	}
	testDeepEqual(t, Point{X: 0.0, Y: 1.0, Z: -1.0})

	type LinkedListNode struct {
		Value Point
		Next  *LinkedListNode
	}
	l := &LinkedListNode{
		Value: Point{X: 0.0, Y: 1.0, Z: -1.0},
		Next:  nil,
	}
	l = &LinkedListNode{
		Value: Point{X: 10, Y: 11, Z: 12},
		Next:  l,
	}
	testDeepEqual(t, l)
}

func TestOf_supports_cyclic_structs_without_unexported_fields(t *testing.T) {
	type LinkedListNode struct {
		Value int
		Next  *LinkedListNode
	}
	first := &LinkedListNode{
		Value: 0,
		Next:  nil,
	}
	second := &LinkedListNode{
		Value: 1,
		Next:  first,
	}
	first.Next = second

	c := clone.Of(first)
	assert.Same(t, c, c.Next.Next)
}

func TestOf_supports_immutable_structs_with_unexported_fields(t *testing.T) {
	// Unexported fields would panic if the type wasn't assumed immutable.
	type ImmutableType struct {
		unexportedField int
	}
	v := ImmutableType{unexportedField: 42}
	clone.AsImmutableType(ImmutableType{})

	c := clone.Of(v)
	assert.Equal(t, v, c)
}

func TestOf_supports_structs_with_unexported_zero_size_traits(t *testing.T) {
	// Unexported fields would panic if the type wasn't assumed immutable.
	type Type struct {
		_     typetraits.NoCompare // Zero-sized unexported field.
		Value int
	}
	v := Type{Value: 42}

	c := clone.Of(v)
	assert.Equal(t, 42, c.Value)
}

func TestOf_immutable_struct_types_are_shallow_copied(t *testing.T) {
	type Immutable struct {
		Slice   []int
		Map     map[int]struct{}
		Pointer *int
	}
	clone.AsImmutableType(Immutable{})

	v := Immutable{
		Slice:   []int{0, 1, 2},
		Map:     map[int]struct{}{0: {}, 1: {}, 2: {}},
		Pointer: new(int),
	}
	c := clone.Of(v)
	assert.DeepEqual(t, v, c)
	assert.Same(t, v.Slice, c.Slice)
	assert.Same(t, v.Map, c.Map)
	assert.Same(t, v.Pointer, c.Pointer)
}

func TestOf_mutable_struct_types_are_deep_copied(t *testing.T) {
	type Mutable struct {
		Slice   []int
		Map     map[int]struct{}
		Pointer *int
	}

	v := Mutable{
		Slice:   []int{0, 1, 2},
		Map:     map[int]struct{}{0: {}, 1: {}, 2: {}},
		Pointer: new(int),
	}
	c := clone.Of(v)
	assert.DeepEqual(t, v, c)
	assert.NotSame(t, v.Slice, c.Slice)
	assert.NotSame(t, v.Map, c.Map)
	assert.NotSame(t, v.Pointer, c.Pointer)
}

func TestOf_panics_on_channels(t *testing.T) {
	defer handleUnsupportedKind(t, reflect.Chan)
	v := chan struct{}(nil)
	clone.Of(v)
}

func TestOf_panics_on_functions(t *testing.T) {
	defer handleUnsupportedKind(t, reflect.Func)
	v := func() {}
	clone.Of(v)
}

func TestOf_panics_on_unsafe_pointers(t *testing.T) {
	defer handleUnsupportedKind(t, reflect.UnsafePointer)
	v := unsafe.Pointer(nil)
	clone.Of(v)
}

func handleUnsupportedKind(t *testing.T, kind reflect.Kind) {
	t.Helper()

	r := recover()
	require.NotNil(t, nil)
	s := require.As[string](t, r)

	const reasonForFailure = "unsupported kind"
	assert.Substring(t, s, reasonForFailure)
	assert.Substring(t, s, kind.String())
}

func TestOf_panics_on_mutable_struct_with_unexported_fields(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		s := require.As[string](t, r)

		const why = "struct has unexported fields"
		assert.Substring(t, s, why)

		const structName = "StructWithPrivateFields"
		assert.Substring(t, s, structName)

		const unexportedFieldName = "anUnexportedField"
		assert.Substring(t, s, unexportedFieldName)
	}()

	type StructWithPrivateFields struct {
		Public            int
		anUnexportedField int
	}
	clone.Of(StructWithPrivateFields{Public: 1, anUnexportedField: 2})
}

func testDeepEqual(t *testing.T, v any) {
	t.Helper()
	t.Run(fmt.Sprintf("%T(%v)", v, v), func(t *testing.T) {
		c := clone.Of(v)
		assert.DeepEqual(t, v, c)
	})
}
