package apperror_test

import (
	"artk.dev/core/apperror"
	"strings"
	"testing"
)

func TestKindValues(t *testing.T) {
	const expected = 10
	got := len(apperror.KindValues())
	if expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestKindOf_returns_OK_for_nil(t *testing.T) {
	got := apperror.KindOf(nil)
	if expected := apperror.OK; expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestKind_String_contains_unexpected_values(t *testing.T) {
	if !strings.Contains(apperror.Kind(-1).String(), "-1") {
		t.Error("unexpected value (-1) not included")
	}
}
