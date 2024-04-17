package apperror_test

import (
	"artk.dev/core/apperror"
	"testing"
)

func TestConflict(t *testing.T) {
	err := apperror.Conflict("%v error", "conflict")
	if k := apperror.KindOf(err); k != apperror.ConflictKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsConflict(err) {
		t.Errorf("expected conflict error, got %v", err)
	}
	if msg := err.Error(); msg != "conflict error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
