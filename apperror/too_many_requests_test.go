package apperror_test

import (
	"github.com/jespert/artk/apperror"
	"testing"
)

func TestTooManyRequests(t *testing.T) {
	err := apperror.TooManyRequests("%v error", "test")
	if k := apperror.KindOf(err); k != apperror.TooManyRequestsKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsTooManyRequests(err) {
		t.Errorf("expected tooManyRequests error, got %v", err)
	}
	if msg := err.Error(); msg != "test error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
