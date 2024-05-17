//go:build race

package racechecker_test

import (
	"artk.dev/racechecker"
	"testing"
)

func TestRequire(t *testing.T) {
	var fakeT T
	racechecker.Require(&fakeT)

	if n := fakeT.NumCallsToHelper; n != 1 {
		t.Errorf("expected 1 call to Helper, got %v", n)
	}
	if n := fakeT.NumCallsToSkip; n != 0 {
		t.Errorf("expected no calls to Skip, got %v", n)
	}
}
