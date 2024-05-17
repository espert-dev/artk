package assume

import "fmt"

// True panics if the condition is false.
func True(ok bool) {
	Truef(ok, "expected condition to be true")
}

// Truef panics if the condition is false.
func Truef(ok bool, format string, args ...any) {
	if !ok {
		panic(fmt.Sprintf(format, args...))
	}
}
