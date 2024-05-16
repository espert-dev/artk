package asynctx_test

import (
	"artk.dev/asynctx"
	"context"
	"testing"
	"time"
)

func TestFrom_Deadline_returns_zero_values(t *testing.T) {
	const timeout = 24 * time.Hour
	parent, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	derived := asynctx.From(parent)
	deadline, ok := derived.Deadline()
	if !deadline.IsZero() {
		t.Error("expected the deadline to be zero, got", deadline)
	}
	if ok {
		t.Error("expected ok to be false, got true")
	}
}

func TestFrom_Done_returns_nil(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()

	derived := asynctx.From(parent)
	select {
	case <-derived.Done():
		t.Error("unexpectedly got a value from Done")
	default:
		// Test succeeds.
	}
}

func TestFrom_Err_returns_nil(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()

	derived := asynctx.From(parent)
	if err := derived.Err(); err != nil {
		t.Error("got unexpected error:", err)
	}
}

func TestFrom_Value_same_as_that_of_parent(t *testing.T) {
	for _, tc := range []struct {
		name string
		kvs  map[any]any
	}{
		{
			name: "none",
			kvs:  nil,
		},
		{
			name: "one",
			kvs: map[any]any{
				"foo": "bar",
			},
		},
		{
			name: "two",
			kvs: map[any]any{
				"foo": "bar",
				"bar": "baz",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Create parent context.
			parent := context.Background()
			for k, v := range tc.kvs {
				parent = context.WithValue(parent, k, v)
			}

			// Check that all values in the derived context match.
			derived := asynctx.From(parent)
			for key, expected := range tc.kvs {
				got := derived.Value(key)
				if expected != got {
					t.Errorf(
						"key %v: expected %v, got %v",
						key,
						expected,
						got,
					)
				}
			}
		})
	}
}
