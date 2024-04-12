package apperror_test

import (
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestValidation(t *testing.T) {
	err := apperror.Validation("test error")
	if !apperror.IsValidation(err) {
		t.Errorf("expected forbidden error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}