package event_test

import (
	"artk.dev/event"
	"artk.dev/racechecker"
	"artk.dev/testbarrier"
	"context"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestStream_all_observers_receive_the_event(t *testing.T) {
	t.Parallel()

	run := func(t *testing.T, n int) {
		t.Parallel()

		t.Logf("Given there are %v registered observers,", n)
		events := make([]Event, n)
		wg := &sync.WaitGroup{}
		wg.Add(n)

		stream := &event.Stream[Event]{}
		for i := range n {
			stream.WillNotify(func(
				_ context.Context,
				e Event,
			) error {
				defer wg.Done()
				events[i] = e
				return nil
			})
		}

		t.Log("When an event is observed,")
		err := stream.Observe(context.TODO(), exampleEvent())
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

func TestStream_observer_errors_are_not_returned_to_caller(t *testing.T) {
	t.Parallel()

	t.Log("When an the observer fails,")
	stream := &event.Stream[Event]{}
	stream.WillNotify(func(_ context.Context, _ Event) error {
		return errors.New("expected observer failure")
	})

	t.Log("Stream will not propagate the error.")
	err := stream.Observe(context.TODO(), exampleEvent())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStream_WillNotify_and_Observe_are_thread_safe(t *testing.T) {
	t.Parallel()
	racechecker.Require(t)

	stream := &event.Stream[Event]{}
	for range 100 {
		go func() {
			stream.WillNotify(func(
				_ context.Context,
				_ Event,
			) error {
				// The observer is irrelevant, only the
				// internal state of the Stream matters.
				return nil
			})
		}()
		go func() {
			err := stream.Observe(context.TODO(), exampleEvent())
			if err != nil {
				t.Error("Unexpected error:", err)
			}
		}()
	}
}

func TestStream_Observe_deep_copies_events(t *testing.T) {
	t.Parallel()

	t.Run("slices", func(t *testing.T) {
		t.Parallel()
		t.Helper()
		assertIsDeepCopy(t, []int{1, 2, 3})
	})
	t.Run("maps", func(t *testing.T) {
		t.Parallel()
		t.Helper()
		assertIsDeepCopy(t, map[string]int{
			"foo": 1,
			"bar": 2,
			"baz": 3,
		})
	})
}

func TestStream_Shutdown_allows_all_tasks_to_terminate(t *testing.T) {
	t.Parallel()

	stream := &event.Stream[Event]{}

	t.Log("Given there are 100 registered observers,")
	const numObservers = 100
	var numFinished atomic.Int64
	for range numObservers {
		stream.WillNotify(func(_ context.Context, _ Event) error {
			numFinished.Add(1)
			return nil
		})
	}

	t.Log("And the Stream has observed an event,")
	err := stream.Observe(context.TODO(), exampleEvent())
	if err != nil {
		t.Error("unexpected error:", err)
	}

	t.Log("When the Stream is shut down,")
	var wg sync.WaitGroup
	stream.Shutdown(&wg)
	wg.Wait()

	t.Log("All 100 registered observers finish normally.")
	if n := numFinished.Load(); n != numObservers {
		t.Errorf("expected %v, got %v", numObservers, n)
	}
}

func TestStream_observer_contexts_inherit_values(t *testing.T) {
	t.Parallel()

	const value = 42
	ctx := context.WithValue(context.TODO(), key{}, value)

	var wg sync.WaitGroup
	stream := &event.Stream[Event]{}

	for range 100 {
		stream.WillNotify(func(ctx context.Context, _ Event) error {
			if ctx.Value(key{}) != 42 {
				t.Errorf("context value not inherited")
			}
			return nil
		})
	}

	if err := stream.Observe(ctx, exampleEvent()); err != nil {
		t.Error("unexpected error:", err)
	}

	stream.Shutdown(&wg)
	wg.Wait()
}

func TestStream_observer_contexts_do_not_inherit_deadline(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	stream := &event.Stream[Event]{}

	for range 100 {
		stream.WillNotify(func(ctx context.Context, _ Event) error {
			if deadline, ok := ctx.Deadline(); ok {
				t.Error("unexpected deadline:", deadline)
			}
			if err := ctx.Err(); err != nil {
				t.Error("unexpected error:", err)
			}
			return nil
		})
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 24*time.Hour)
	defer cancel()

	if err := stream.Observe(ctx, exampleEvent()); err != nil {
		t.Error("unexpected error:", err)
	}

	stream.Shutdown(&wg)
	wg.Wait()
}

func TestStream_Observe_does_not_notify_if_context_is_cancelled(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	stream := &event.Stream[Event]{}

	for range 100 {
		stream.WillNotify(func(_ context.Context, _ Event) error {
			t.Error("observer incorrectly notified")
			return nil
		})
	}

	ctx, cancel := context.WithCancel(context.TODO())

	// The context is cancelled from the start.
	cancel()

	if err := stream.Observe(ctx, exampleEvent()); err != nil {
		t.Error("unexpected error:", err)
	}

	stream.Shutdown(&wg)
	wg.Wait()
}

func TestStream_supports_context_middleware(t *testing.T) {
	t.Parallel()

	var key key
	const expected = 42

	t.Log("Given that the middleware will insert a known key-value,")
	contextMiddleware := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key, expected)
	}

	t.Log("Then the observer will receive the known key-value")
	barrier := testbarrier.New()
	stream := &event.Stream[Event]{}
	stream.WithContextMiddleware(contextMiddleware)
	stream.WillNotify(func(ctx context.Context, _ Event) error {
		defer barrier.Lift()
		got, ok := ctx.Value(key).(int)
		if !ok {
			t.Error("expected key-value not found")
		}
		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
		return nil
	})

	t.Log("When Stream observes an error.")
	err := stream.Observe(context.TODO(), exampleEvent())
	if err != nil {
		t.Error("unexpected error:", err)
	}

	barrier.Wait(t, 5*time.Second)
}

func TestStream_observer_middleware_can_modify_context(t *testing.T) {
	t.Parallel()

	var key key
	const expected = 42

	t.Log("Given that the middleware will insert a known key-value,")
	observerMiddleware := func(
		next event.Observer[Event],
	) event.Observer[Event] {
		return func(ctx context.Context, e Event) error {
			return next(context.WithValue(ctx, key, expected), e)
		}
	}

	t.Log("Then the observer will receive the known key-value")
	barrier := testbarrier.New()
	stream := &event.Stream[Event]{}
	stream.WithObserverMiddleware(observerMiddleware)
	stream.WillNotify(func(ctx context.Context, _ Event) error {
		defer barrier.Lift()
		got, ok := ctx.Value(key).(int)
		if !ok {
			t.Error("expected key-value not found")
		}
		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
		return nil
	})

	t.Log("When Stream observes an error.")
	err := stream.Observe(context.TODO(), exampleEvent())
	if err != nil {
		t.Error("unexpected error:", err)
	}

	barrier.Wait(t, 5*time.Second)
}

func TestStream_context_middleware_happens_before_observer(t *testing.T) {
	t.Parallel()

	var key key
	const original = 13
	const expected = 42

	t.Logf("Given that the context middleware will insert %v", original)
	contextMiddleware := func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key, original)
	}

	t.Logf("And that the observer middleware will insert %v,", expected)
	observerMiddleware := func(
		next event.Observer[Event],
	) event.Observer[Event] {
		return func(ctx context.Context, e Event) error {
			// Overwrites value set by the context middleware.
			return next(context.WithValue(ctx, key, expected), e)
		}
	}

	t.Logf("Then the observer will receive the value %v", expected)
	barrier := testbarrier.New()
	stream := &event.Stream[Event]{}
	stream.WithContextMiddleware(contextMiddleware)
	stream.WithObserverMiddleware(observerMiddleware)
	stream.WillNotify(func(ctx context.Context, _ Event) error {
		defer barrier.Lift()
		got, ok := ctx.Value(key).(int)
		if !ok {
			t.Error("expected key-value not found")
		}
		if got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
		return nil
	})

	t.Log("When Stream observes an error.")
	err := stream.Observe(context.TODO(), exampleEvent())
	if err != nil {
		t.Error("unexpected error:", err)
	}

	barrier.Wait(t, 5*time.Second)
}

func TestStream_messages_are_dropped_when_queue_size_exceeded(t *testing.T) {
	const queueSize = 0
	stream := event.NewStream[Event](queueSize)

	var numEventsObserved int
	barrier := make(chan struct{})
	stream.WillNotify(func(_ context.Context, _ Event) error {
		<-barrier
		numEventsObserved++
		return nil
	})

	// Just give a big enough number that we can be confident about
	// the message dropping behaviour.
	const excessAmount = 100
	for i := range queueSize + excessAmount {
		e := Event{
			ID:   i,
			Name: "test",
		}
		if err := stream.Observe(context.TODO(), e); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// The barrier is not released all events have been notified.
	close(barrier)

	// Shutdown the broker to ensure exhaustive processing.
	var wg sync.WaitGroup
	stream.Shutdown(&wg)
	wg.Wait()

	// The special behaviour of a channel of size zero means that we can
	// only guarantee that the observed events will be either queueSize
	// or (queueSize+1).
	if numEventsObserved != queueSize && numEventsObserved != queueSize+1 {
		t.Errorf(
			"expected %v or %v, got %v",
			queueSize,
			queueSize+1,
			numEventsObserved,
		)
	}
}
