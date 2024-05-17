package assume_test

import "testing"

func expectPanic(t *testing.T, expected string, fn func()) {
	t.Helper()

	defer func(t *testing.T) {
		t.Helper()

		r := recover()
		if r == nil {
			t.Error("Missing expected panic")
		}

		got, ok := r.(string)
		if !ok {
			t.Error("Panic object is not a string")
		}

		if expected != got {
			t.Errorf(
				`Expected message "%v", got "%v"`,
				expected,
				got,
			)
		}
	}(t)

	fn()
}

func expectNoPanic(t *testing.T, fn func()) {
	t.Helper()

	defer func(t *testing.T) {
		t.Helper()
		if r := recover(); r != nil {
			t.Error("Unexpected panic")
		}
	}(t)

	fn()
}

// Strings for custom formatted messages.
const (
	expectedCustomMessage = "value: 42"
	format                = "value: %v"
	value                 = 42
)
