package clone_test

import (
	"fmt"
	"github.com/jespert/artk/clone"
	"math"
	"reflect"
	"testing"
)

// Derived types to ensure that cloning works across the kind.
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

func TestOf_trivial_copy(t *testing.T) {
	// Don't test with NaN because comparison is always false.
	pInf := math.Inf(1)
	nInf := math.Inf(-1)

	testValues(t, []bool{false, true})
	testValues(t, []boolType{false, true})
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
				if got := clone.Of(v); got != v {
					t.Errorf(
						"expected %v, got %v",
						v,
						got,
					)
				}
			})
		}
	})
}

func TestOf_arrays(t *testing.T) {
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
		if !reflect.DeepEqual(array, c) {
			t.Errorf("the arrays are different")
		}
	})
}

func TestOf_interfaces(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var v any
		c := clone.Of(v)
		if v != c {
			t.Errorf("expected %v, got %v", v, c)
		}
	})
	t.Run("not nil", func(t *testing.T) {
		var v any = true
		c := clone.Of(v)
		if v != c {
			t.Errorf("expected %v, got %v", v, c)
		}
	})
}

func TestOf_pointers(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var v *int

		c := clone.Of(v)
		if c != nil {
			t.Errorf("unexpected not nil")
		}
	})
	t.Run("not nil", func(t *testing.T) {
		v := new(int)
		*v = 1

		c := clone.Of(v)
		if c == v {
			t.Errorf("unexpected shallow copy")
		}
		if *c != *v {
			t.Errorf("expected %v, got %v", *v, *c)
		}
	})
}

func TestOf_slices(t *testing.T) {
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
	var name string
	if slice == nil {
		name = fmt.Sprintf("%T(nil)", slice)
	} else {
		name = fmt.Sprintf("%T%v", slice, slice)
	}

	t.Helper()
	t.Run(name, func(t *testing.T) {
		c := clone.Of(slice)

		if slice == nil && c != nil || slice != nil && c == nil {
			t.Errorf("nil not preserved")
		}
		if len(c) != len(slice) {
			t.Errorf("the slices have different lengths")
		}
		if !reflect.DeepEqual(slice, c) {
			t.Errorf("the slices have different elements")
		}

		// Sharing is never an issue if the slices are empty.
		if len(slice) != 0 {
			ps := reflect.ValueOf(slice).Pointer()
			pc := reflect.ValueOf(c).Pointer()
			if ps == pc {
				t.Errorf("the slices share memory")
			}
		}
	})
}

func TestOf_string(t *testing.T) {
	testString(t, "")
	testString(t, stringType(""))
	testString(t, "foo")
	testString(t, stringType("foo"))
}

func testString[T ~string](t *testing.T, s T) {
	t.Helper()
	t.Run(fmt.Sprintf("%T(%v)", s, s), func(t *testing.T) {
		if c := clone.Of(s); c != s {
			t.Errorf("expected %v, got %v", s, c)
		}
	})
}
