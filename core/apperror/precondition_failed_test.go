package apperror_test

import (
	"artk.dev/core/apperror"
	"testing"
)

func TestPreconditionFailed(t *testing.T) {
	err := apperror.PreconditionFailed("%v error", "precondition failed")
	if k := apperror.KindOf(err); k != apperror.PreconditionFailedKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsPreconditionFailed(err) {
		t.Errorf("expected preconditionFailed error, got %v", err)
	}
	if msg := err.Error(); msg != "precondition failed error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}
