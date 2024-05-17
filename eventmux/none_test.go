package eventmux_test

import (
	"artk.dev/eventmux"
	"context"
	"testing"
)

func TestNone(t *testing.T) {
	err := eventmux.None[Event](context.TODO(), exampleEvent())
	if err != nil {
		t.Error("unexpected error:", err)
	}
}
