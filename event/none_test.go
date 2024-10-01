package event_test

import (
	"artk.dev/event"
	"context"
	"testing"
)

func TestNone(t *testing.T) {
	err := event.None[Event](context.TODO(), exampleEvent())
	if err != nil {
		t.Error("unexpected error:", err)
	}
}
