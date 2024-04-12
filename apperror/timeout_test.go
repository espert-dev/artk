package apperror_test

import (
	"context"
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestTimeout(t *testing.T) {
	err := apperror.Timeout("test error")
	if !apperror.IsTimeout(err) {
		t.Errorf("expected timeout error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}

func TestContextDeadlineExceededIsTimeout(t *testing.T) {
	if !apperror.IsTimeout(context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded to be a timeout")
	}
}
