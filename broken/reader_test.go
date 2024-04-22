package broken_test

import (
	"artk.dev/broken"
	"testing"
)

func TestReader_Read_always_fails(t *testing.T) {
	var buffer []byte
	var r broken.Reader
	n, err := r.Read(buffer)
	if n != 0 {
		t.Errorf("expected 0, got %v", n)
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}
}
