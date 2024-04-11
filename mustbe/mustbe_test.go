package mustbe_test

import (
	"fmt"
	"github.com/jespert/artk/mustbe"
	"testing"
)

func TestEqual(t *testing.T) {
	t.Run("Do not panic if equal", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Unexpected panic")
			}
		}()

		mustbe.Equal(true, true)
	})
	t.Run("Panic if different", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Missing expected panic")
			}
		}()

		mustbe.Equal(false, true)
	})
}

func TestNoError(t *testing.T) {
	t.Run("Do not panic if nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Unexpected panic")
			}
		}()

		mustbe.NoError(nil)
	})
	t.Run("Panic if not nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Missing expected panic")
			}
		}()

		mustbe.NoError(fmt.Errorf("test error"))
	})
}

func TestNotNil(t *testing.T) {
	t.Run("Do not panic if not nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Unexpected panic")
			}
		}()

		mustbe.NotNil(struct{}{})
	})
	t.Run("Panic if nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Missing expected panic")
			}
		}()

		mustbe.NotNil(nil)
	})
}

func TestTrue(t *testing.T) {
	t.Run("Do not panic if true", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Unexpected panic")
			}
		}()

		mustbe.True(true)
	})
	t.Run("Panic if false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Missing expected panic")
			}
		}()

		mustbe.True(false)
	})
}
