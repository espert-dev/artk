package apperror_test

import (
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestConflict(t *testing.T) {
	err := apperror.Conflict("test error")
	if k := apperror.KindOf(err); k != apperror.ConflictKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsConflict(err) {
		t.Errorf("expected conflict error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
