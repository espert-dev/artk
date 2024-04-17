package apperror_test

import (
	"artk.dev/core/apperror"
	"context"
	"testing"
)

func TestTimeout(t *testing.T) {
	err := apperror.Timeout("%v error", "timeout")
	if k := apperror.KindOf(err); k != apperror.TimeoutKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsTimeout(err) {
		t.Errorf("expected timeout error, got %v", err)
	}
	if msg := err.Error(); msg != "timeout error" {
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
