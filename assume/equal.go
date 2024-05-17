package assume

import "fmt"

// Equal panics if the items are not equal.
func Equal[T comparable](x, y T) {
	Equalf(x, y, "expected %v == %v", x, y)
}

// Equalf panics if the items are not equal.
func Equalf[T comparable](x, y T, format string, args ...any) {
	if x != y {
		panic(fmt.Sprintf(format, args...))
	}
}
