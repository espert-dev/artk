package assume_test

import (
	"artk.dev/assume"
	"testing"
)

func TestNotZero(t *testing.T) {
	t.Run("Do not panic if not zero", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.NotZero(new(int))
		})
	})
	t.Run("Panic if zero", func(t *testing.T) {
		expectPanic(t, "zero value", func() {
			assume.NotZero((*int)(nil))
		})
	})
}

func TestNotZerof(t *testing.T) {
	t.Run("Do not panic if not zero", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.NotZerof(new(int), format, value)
		})
	})
	t.Run("Panic if zero", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assume.NotZerof((*int)(nil), format, value)
		})
	})
}

func TestNotNilSlice(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.NotNilSlice([]int{})
		})
	})
	t.Run("Panic if nil", func(t *testing.T) {
		expectPanic(t, "nil slice", func() {
			assume.NotNilSlice([]int(nil))
		})
	})
}

func TestNotNilSlicef(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.NotNilSlicef([]int{}, format, value)
		})
	})
	t.Run("Panic if nil", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assume.NotNilSlicef([]int(nil), format, value)
		})
	})
}

func TestNotNilMap(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.NotNilMap(map[bool]int{})
		})
	})
	t.Run("Panic if zero", func(t *testing.T) {
		expectPanic(t, "nil map", func() {
			assume.NotNilMap(map[bool]int(nil))
		})
	})
}

func TestNotNilMapf(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.NotNilMapf(map[bool]int{}, format, value)
		})
	})
	t.Run("Panic if nil", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assume.NotNilMapf(map[bool]int(nil), format, value)
		})
	})
}
