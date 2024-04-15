package apperror_test

import (
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestUnauthorized(t *testing.T) {
	err := apperror.Unauthorized("%v error", "test")
	if k := apperror.KindOf(err); k != apperror.UnauthorizedKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsUnauthorized(err) {
		t.Errorf("expected unauthorized error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
