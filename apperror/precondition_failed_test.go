package apperror_test

import (
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestPreconditionFailed(t *testing.T) {
	err := apperror.PreconditionFailed("test error")
	if !apperror.IsPreconditionFailed(err) {
		t.Errorf("expected preconditionFailed error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}