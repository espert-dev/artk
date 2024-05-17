package assume_test

import (
	"artk.dev/assume"
	"errors"
	"testing"
)

func TestSuccess(t *testing.T) {
	t.Run("Do not panic if nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.Success(nil)
		})
	})
	t.Run("Panic if not nil", func(t *testing.T) {
		expectPanic(t, "unexpected error: test error", func() {
			assume.Success(errors.New("test error"))
		})
	})
}

func TestSuccessf(t *testing.T) {
	t.Run("Do not panic if nil", func(t *testing.T) {
		expectNoPanic(t, func() {
			assume.Successf(nil, format, value)
		})
	})
	t.Run("Panic if not nil", func(t *testing.T) {
		expectPanic(t, expectedCustomMessage, func() {
			err := errors.New("test error")
			assume.Successf(err, format, value)
		})
	})
}
