package apperror_test

import (
	"context"
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestTimeout(t *testing.T) {
	err := apperror.Timeout("%v error", "test")
	if k := apperror.KindOf(err); k != apperror.TimeoutKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsTimeout(err) {
		t.Errorf("expected timeout error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}

func TestContextDeadlineExceededIsTimeout(t *testing.T) {
	err := context.DeadlineExceeded
	if k := apperror.KindOf(err); k != apperror.TimeoutKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsTimeout(err) {
		t.Errorf("expected context.DeadlineExceeded to be a timeout")
	}
}
