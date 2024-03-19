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
	t.Helper()

	var name string
	if slice == nil {
		name = fmt.Sprintf("%T(nil)", slice)
	} else {
		name = fmt.Sprintf("%T%v", slice, slice)
	}

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
	})
}
