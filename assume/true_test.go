package assume_test

import (
	"artk.dev/assume"
	"testing"
)

func TestTrue(t *testing.T) {
	t.Run("Do not panic if true", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.True(true)
		})
	})
	t.Run("Panic if false", func(t *testing.T) {
		expectPanic(t, "expected condition to be true", func() {
			assume.True(false)
		})
	})
}

func TestTruef(t *testing.T) {
	t.Run("Do not panic if true", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.Truef(true, format, value)
		})
	})
	t.Run("Panic if false", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			assume.Truef(false, format, value)
		})
	})
}
