// Package racechecker can be used to skip tests that require the race
// condition checker if it is not enabled in the current build.
package racechecker

type testingT interface {
	Skip(args ...any)
}
