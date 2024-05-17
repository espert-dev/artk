package eventmux

import (
	"artk.dev/clone"
	"context"
	"sync"
)

var _ Observer[any] = (&Mux[any]{}).Observe

// Mux is a thread-safe in-memory event multiplexer.
type Mux[Event any] struct {
	mutex sync.RWMutex   // 24 bytes on 64 bits.
	wg    sync.WaitGroup // 12 bytes on 64 bits.

	// On 64-bit systems, 4 bytes of padding will be inserted here to
	// ensure that 64-bit words remain aligned.

	observers []Observer[Event] // 24 bytes on 64 bits.
}

// Observe and propagate an event to registered observers.
func (m *Mux[Event]) Observe(ctx context.Context, event Event) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	numObservers := len(m.observers)
	m.wg.Add(numObservers)

	// Observers are notified concurrently.
	for i := range numObservers {
		// Trade performance for safety.
		event := clone.Of(event)

		go func(
			ctx context.Context,
			observer Observer[Event],
			event Event,
		) {
			defer m.wg.Done()

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

// Shutdown the Mux and communicate finishing via the sync.WaitGroup.
func (m *Mux[Event]) Shutdown(wg *sync.WaitGroup) {
	// Synchronously prevent new messages from being sent.
	m.stopEventPropagation()

	// Asynchronously wait for tasks to finish.
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.wg.Wait()
	}()
}

func (m *Mux[Event]) stopEventPropagation() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.observers = nil
}

// New creates a Mux.
func New[Event any]() *Mux[Event] {
	return &Mux[Event]{}
}
