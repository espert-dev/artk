package apperror_test

import (
	"artk.dev/apperror"
	"testing"
)

func TestForbidden(t *testing.T) {
	err := apperror.Forbidden("%v error", "test")
	if k := apperror.KindOf(err); k != apperror.ForbiddenKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsForbidden(err) {
		t.Errorf("expected forbidden error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
