package typetraits_test

import (
	"artk.dev/typetraits"
	"testing"
	"unsafe"
)

func TestNoCopy_has_size_zero(t *testing.T) {
	if size := unsafe.Sizeof(typetraits.NoCopy{}); size != 0 {
		t.Errorf("expected 0, got %v", size)
	}
}

func TestNoCompare_has_size_zero(t *testing.T) {
	if size := unsafe.Sizeof(typetraits.NoCompare{}); size != 0 {
		t.Errorf("expected 0, got %v", size)
	}
}
