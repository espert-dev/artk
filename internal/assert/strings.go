package assert

import "strings"

func Substring[S ~string](t T, s, substr S) bool {
	t.Helper()

	if !strings.Contains(string(s), string(substr)) {
		report2(t, "substring not found", substr, s)
		return false
	}

	return true
}

func NotSubstring[S ~string](t T, s, substr S) bool {
	t.Helper()

	if strings.Contains(string(s), string(substr)) {
		report2(t, "substring found", substr, s)
		return false
	}

	return true
}
