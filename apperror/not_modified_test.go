package apperror_test

import (
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestNotModified(t *testing.T) {
	err := apperror.NotModified("test error")
	if !apperror.IsNotModified(err) {
		t.Errorf("expected notModified error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
