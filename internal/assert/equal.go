package assert

import (
	"fmt"
	"reflect"
)

func Equal[K comparable](t T, expected, got K) bool {
	t.Helper()

	if expected != got {
		report2(t, "not equal", expected, got)
		return false
	}

	return true
}

func NotEqual[K comparable](t T, x, y K) bool {
	t.Helper()

	if x == y {
		report1(t, "equal", x)
		return false
	}

	return true
}

func DeepEqual[V any](t T, expected, got V) bool {
	t.Helper()

	if !reflect.DeepEqual(expected, got) {
		report2(t, "not deep equal", expected, got)
		return false
	}

	return true
}

func NotDeepEqual[V any](t T, x, y V) bool {
	t.Helper()

	if reflect.DeepEqual(x, y) {
		report1(t, "deep equal", x)
		return false
	}

	return true
}

func Same[V any](t T, expected, got V) bool {
	t.Helper()

	if !same(expected, got) {
		report2(
			t,
			"not same",
			fmt.Sprintf("%v", expected), // TODO use pointer
			fmt.Sprintf("%v", got),      // TODO use pointer
		)
		return false
	}

	return true
}

func NotSame[V any](t T, x, y V) bool {
	t.Helper()

	if same(x, y) {
		report1(
			t,
			"same",
			fmt.Sprintf("%v", x), // TODO use pointer
		)
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
