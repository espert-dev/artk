package assume

import "fmt"

// Success panics if the provided error is not nil.
func Success(err error) {
	Successf(err, "unexpected error: %v", err)
}

// Successf panics if the provided error is not nil.
func Successf(err error, format string, args ...any) {
	if err != nil {
		panic(fmt.Sprintf(format, args...))
	}
}
