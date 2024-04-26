package assert

import "fmt"

// Is attempts a cast of the value x to the type V.
// If it succeeds, it returns the cast value and true.
// Otherwise, it returns the zero value for V and false.
func Is[V any](t T, x any) (V, bool) {
	t.Helper()

	v, ok := x.(V)
	if !ok {
		var zero V
		report2(
			t,
			"type is not convertible",
			fmt.Sprintf("%T", zero),
			fmt.Sprintf("%T", x),
		)
		return zero, false
	}

	return v, true
}
