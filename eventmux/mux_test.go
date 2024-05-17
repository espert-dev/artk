package eventmux_test

import (
	"artk.dev/eventmux"
	"artk.dev/racechecker"
	"artk.dev/testbarrier"
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Event struct {
	ID   int
	Name string
}

func TestMux_all_observers_receive_the_event(t *testing.T) {
	t.Parallel()

	run := func(t *testing.T, n int) {
		t.Parallel()

		t.Logf("Given there are %v registered observers,", n)
		events := make([]Event, n)
		wg := &sync.WaitGroup{}
		wg.Add(n)

		mux := eventmux.New[Event]()
		for i := range n {
			mux.WillNotify(func(_ context.Context, e Event) error {
				defer wg.Done()
				events[i] = e
				return nil
			})
		}

		t.Log("When an event is observed,")
		err := mux.Observe(context.TODO(), exampleEvent())
		if err != nil {
			t.Fatal("unexpected Observe error")
		}

		t.Log("Then, eventually, all observers will be notified")
		testbarrier.WaitForGroup(t, wg, 5*time.Second)

		t.Log("And they will have received the expected event.")
		expected := exampleEvent()
		for i := range n {
			if expected != events[i] {
				t.Errorf(
					"[%v] expected %v, got %v",
					exampleEvent(),
					expected,
					events[i],
				)
			}
		}
	}
	for n := range 10 {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			t.Helper()
			run(t, n)
		})
	}
}

func TestMux_observer_errors_are_not_returned_to_caller(t *testing.T) {
	t.Parallel()

	t.Log("When an the observer fails,")
	mux := eventmux.New[Event]()
	mux.WillNotify(func(_ context.Context, _ Event) error {
		return errors.New("expected observer failure")
	})

	t.Log("Mux will not propagate the error.")
	err := mux.Observe(context.TODO(), exampleEvent())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMux_WillNotify_and_Observe_are_thread_safe(t *testing.T) {
	t.Parallel()
	racechecker.Require(t)

	mux := eventmux.New[Event]()
	for range 100 {
		go func() {
			mux.WillNotify(func(_ context.Context, _ Event) error {
				// The observer is irrelevant, only the
				// internal state of the Mux matters.
				return nil
			})
		}()
		go func() {
			err := mux.Observe(context.TODO(), exampleEvent())
			if err != nil {
				t.Error("Unexpected error:", err)
			}
		}()
	}
}

func exampleEvent() Event {
	return Event{
		ID:   expectedID,
		Name: expectedName,
	}
}

const expectedID = 1234
const expectedName = "Test Event"
