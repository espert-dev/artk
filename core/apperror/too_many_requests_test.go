package apperror_test

import (
	"artk.dev/core/apperror"
	"testing"
)

func TestTooManyRequests(t *testing.T) {
	err := apperror.TooManyRequests("%v error", "too many requests")
	if k := apperror.KindOf(err); k != apperror.TooManyRequestsKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsTooManyRequests(err) {
		t.Errorf("expected tooManyRequests error, got %v", err)
	}
	if msg := err.Error(); msg != "too many requests error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
