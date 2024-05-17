package assume_test

import (
	"artk.dev/assume"
	"testing"
)

func TestEqual(t *testing.T) {
	t.Run("Do not panic if equal", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.Equal(true, true)
		})
	})
	t.Run("Panic if different", func(t *testing.T) {
		expectPanic(t, "expected false == true", func() {
			assume.Equal(false, true)
		})
	})
}

func TestEqualf(t *testing.T) {
	t.Run("Do not panic if equal", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.Equalf(true, true, format, value)
		})
	})
	t.Run("Panic if different", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assume.Equalf(false, true, format, value)
		})
	})
}
