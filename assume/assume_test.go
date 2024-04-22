package assume_test

import (
	"artk.dev/assume"
	"errors"
	"testing"
)

func TestEqual(t *testing.T) {
	t.Run("Do not panic if equal", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Unexpected panic")
			}
		}()

		assume.Equal(true, true)
	})
	t.Run("Panic if different", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Missing expected panic")
			}
		}()

		assume.Equal(false, true)
	})
}

func TestSuccess(t *testing.T) {
	t.Run("Do not panic if nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Unexpected panic")
			}
		}()

		assume.Success(nil)
	})
	t.Run("Panic if not nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Missing expected panic")
			}
		}()

		assume.Success(errors.New("test error"))
	})
}

func TestNotNil(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Unexpected panic")
			}
		}()

		assume.NotNil(struct{}{})
	})
	t.Run("Panic if nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Missing expected panic")
			}
		}()

		assume.NotNil(nil)
	})
}

func TestTrue(t *testing.T) {
	t.Run("Do not panic if true", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Unexpected panic")
			}
		}()

		assume.True(true)
	})
	t.Run("Panic if false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Missing expected panic")
			}
		}()

		assume.True(false)
	})
}
