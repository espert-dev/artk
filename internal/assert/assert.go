package assert

import (
	"reflect"
	"strings"
)

// T abstracts the functionality of testing.T used by this package.
type T interface {
	Helper()
	Errorf(format string, args ...any)
}

func As[V any](t T, x any) (V, bool) {
	t.Helper()

	v, ok := x.(V)
	if !ok {
		var zero V
		t.Errorf("expected a %V value, got %V", zero, x)
		return zero, false
	}

	return v, true
}

func Substring[S ~string](t T, s, substr S) bool {
	t.Helper()

	if !strings.Contains(string(s), string(substr)) {
		t.Errorf(`does not contain substring "%v": %v`, substr, s)
		return false
	}

	return true
}

func NotSubstring[S ~string](t T, s, substr S) bool {
	t.Helper()

	if strings.Contains(string(s), string(substr)) {
		t.Errorf(`contains substring "%v": %v`, substr, s)
		return false
	}

	return true
}

func Equal[K comparable](t T, expected, got K) bool {
	t.Helper()

	if expected != got {
		t.Errorf(`expected "%v", got "%v"`, expected, got)
		return false
	}

	return true
}

func NotEqual[K comparable](t T, x, y K) bool {
	t.Helper()

	if x == y {
		t.Errorf(`expected "%v", got "%v"`, x, y)
		return false
	}

	return true
}

func DeepEqual[V any](t T, expected, got V) bool {
	t.Helper()

	if !reflect.DeepEqual(expected, got) {
		t.Errorf(`expected "%v", got "%v"`, expected, got)
		return false
	}

	return true
}

func NotDeepEqual[V any](t T, x, y V) bool {
	t.Helper()

	if reflect.DeepEqual(x, y) {
		t.Errorf(`unexpectedly equal: %v`, x)
		return false
	}

	return true
}

func Same[V any](t T, expected, got V) bool {
	t.Helper()

	if !same(expected, got) {
		t.Errorf(`unexpectedly not sharing memory`)
		return false
	}

	return true
}

func NotSame[V any](t T, x, y V) bool {
	t.Helper()

	if same(x, y) {
		t.Errorf(`unexpectedly sharing memory`)
		return false
	}

	return true
}

// PanicBecause checks that the panic value is a string that matches the reason.
// Must be run deferred.
func PanicBecause(t T, why string) bool {
	t.Helper()

	r := recover()
	if r == nil {
		t.Errorf("missing expected panic")
		return false
	}

	s, ok := r.(string)
	if !ok {
		t.Errorf("expected a string, got: %v", r)
		return false
	}

	if !strings.Contains(s, why) {
		t.Errorf("cannot find expected reason in %v", s)
		return false
	}

	return true
}

// NoPanic checks verifies that no panic happened.
// Meant to be run as deferred.
func NoPanic(t T) bool {
	t.Helper()

	if r := recover(); r != nil {
		t.Errorf("unexpected panic")
		return false
	}

	return true
}

// same checks if the two objects share memory (i.e., are shallow copies).
func same[V any](x, y V) bool {
	vx := reflect.ValueOf(x)
	vy := reflect.ValueOf(y)

	// Checking the validity of the value is necessary for interface types.
	if !vx.IsValid() || !vy.IsValid() {
		return false
	}

	// Checking for nil is necessary for concrete types.
	if vx.IsNil() || vy.IsNil() {
		return false
	}

	return vx.Pointer() == vy.Pointer()
}
