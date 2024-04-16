package apperror_test

import (
	"artk.dev/core/apperror"
	"testing"
)

func TestValidation(t *testing.T) {
	err := apperror.Validation("%v error", "test")
	if k := apperror.KindOf(err); k != apperror.ValidationKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsValidation(err) {
		t.Errorf("expected forbidden error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
