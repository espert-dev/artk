package apperror_test

import (
	"artk.dev/core/apperror"
	"testing"
)

func TestNotFound(t *testing.T) {
	err := apperror.NotFound("%v error", "not found")
	if k := apperror.KindOf(err); k != apperror.NotFoundKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsNotFound(err) {
		t.Errorf("expected forbidden error, got %v", err)
	}
	if msg := err.Error(); msg != "not found error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
