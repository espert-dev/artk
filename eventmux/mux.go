package eventmux

import (
	"artk.dev/clone"
	"context"
	"sync"
)

var _ Observer[any] = (&Mux[any]{}).Observe

// Mux is an in-memory event multiplexer.
type Mux[Event any] struct {
	mutex     sync.RWMutex
	observers []Observer[Event]
}

// Observe and propagate an event to registered observers.
func (m *Mux[Event]) Observe(ctx context.Context, event Event) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Observers are notified concurrently.
	for i := range len(m.observers) {
		// Trade performance for safety.
		event := clone.Of(event)

		go func(
			ctx context.Context,
			observer Observer[Event],
			event Event,
		) {
			// Call the observer.
			//
			// While we ultimately ignore the error here, it was
			// made available to any middleware. This can be used,
			// e.g., for logging.
			_ = observer(ctx, event)
		}(ctx, m.observers[i], event)
	}

	return nil
}

// WillNotify registers an observer.
// All events observed by Mux will be propagated to all registered observers.
func (m *Mux[Event]) WillNotify(
	observers ...Observer[Event],
) *Mux[Event] {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.observers = append(m.observers, observers...)

	// Chaining improves DX.
	return m
}

// New creates a Mux.
func New[Event any]() *Mux[Event] {
	return &Mux[Event]{}
}
