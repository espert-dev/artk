// Package assume provides assertions that panic on violation.
//
// This serves two purposes:
//  1. Failing fast.
//  2. Remove unnecessary branches and the temptation to test them.
package assume

import "fmt"

// Equal panics if the items are not equal.
func Equal[T comparable](x, y T) {
	if x != y {
		panic(fmt.Sprintf(
			"constraint violation: expected %v == %v",
			x,
			y,
		))
	}
}

// Success panics if the provided error is not nil.
func Success(err error) {
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %v", err))
	}
}

// NotNil panics if the value is nil.
func NotNil(v interface{}) {
	if v == nil {
		panic("constraint violation: value cannot be nil")
	}
}

// True panics if the condition is false.
func True(ok bool) {
	if !ok {
		panic("constraint violation: expected condition to be true")
	}
}
