package assert

import "strings"

// PanicBecause checks that the panic value is a string that matches the reason.
// Must be run deferred.
func PanicBecause(t T, why string) bool {
	t.Helper()

	r := recover()
	if r == nil {
		report0(t, "missing expected panic")
		return false
	}

	s, ok := r.(string)
	if !ok {
		report1(t, "expected panic value to be a string", r)
		return false
	}

	if !strings.Contains(s, why) {
		report2(
			t,
			"explanation is not part of panic message",
			why,
			s,
		)
		return false
	}

	return true
}

// NoPanic checks verifies that no panic happened.
// Meant to be run as deferred.
func NoPanic(t T) bool {
	t.Helper()

	if r := recover(); r != nil {
		report0(t, "unexpected panic")
		return false
	}

	return true
}
