// Package racechecker can be used to skip tests that require the race
// condition checker if it is not enabled in the current build.
package racechecker

import (
	"runtime/debug"
)

// Require skips the test if the race condition checker is not enabled.
// It will call t.Error if the build information is not accessible.
//
// It might save some wasted cycles, but it is mainly live documentation.
func Require(t testingT) {
	// Beware: this code is inherently untestable.
	// Be extra careful when you review it.

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		t.Error("Cannot read build information")
	}

	for _, setting := range buildInfo.Settings {
		if setting.Key == "-race" && setting.Value == "true" {
			// The race condition checker is enabled.
			return
		}
	}

	t.Skip("This test requires the race checker. Skipping.")
}

type testingT interface {
	Skip(args ...any)
	Error(args ...any)
}
