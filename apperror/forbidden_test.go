package apperror_test

import (
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestForbidden(t *testing.T) {
	err := apperror.Forbidden("test error")
	if !apperror.IsForbidden(err) {
		t.Errorf("expected forbidden error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
