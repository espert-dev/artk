package apperror_test

import (
	"artk.dev/core/apperror"
	"testing"
)

func TestNotModified(t *testing.T) {
	err := apperror.NotModified("%v error", "not modified")
	if k := apperror.KindOf(err); k != apperror.NotModifiedKind {
		t.Errorf("unexpected kind, got %v", k)
	}
	if !apperror.IsNotModified(err) {
		t.Errorf("expected notModified error, got %v", err)
	}
	if msg := err.Error(); msg != "not modified error" {
		t.Errorf("unexpected error message: %v", msg)
	}
}