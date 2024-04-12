package apperror_test

import (
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestNotFound(t *testing.T) {
	err := apperror.NotFound("test error")
	if !apperror.IsNotFound(err) {
		t.Errorf("expected forbidden error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
