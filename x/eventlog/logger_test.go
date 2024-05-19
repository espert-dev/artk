package eventlog_test

import (
	"artk.dev/x/eventlog"
	"artk.dev/x/testlog"
	"context"
	"errors"
	"testing"
)

// The output test must be manually reviewed on change.
func TestLoggerMiddleware_calls_next_on_success(t *testing.T) {
	var numTimesNextCalled int
	next := func(_ context.Context, _ Event) error {
		numTimesNextCalled++
		return nil
	}

	logger := testlog.New(t)
	observer := eventlog.Logger[Event](logger)(next)
	_ = observer(context.TODO(), exampleEvent())

	if numTimesNextCalled != 1 {
		t.Errorf("expected 1 call, got %v", numTimesNextCalled)
	}
}

// The output test must be manually reviewed on change.
func TestLoggerMiddleware_calls_next_on_failure(t *testing.T) {
	var numTimesNextCalled int
	next := func(_ context.Context, _ Event) error {
		numTimesNextCalled++
		return errors.New("expected test failure")
	}

	logger := testlog.New(t)
	observer := eventlog.Logger[Event](logger)(next)
	_ = observer(context.TODO(), exampleEvent())

	if numTimesNextCalled != 1 {
		t.Errorf("expected 1 call, got %v", numTimesNextCalled)
	}
}

type Event struct {
	Key   string
	Value int
}

func exampleEvent() Event {
	return Event{
		Key:   "foo",
		Value: 42,
	}
}
