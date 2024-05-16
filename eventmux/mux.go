package eventmux

import (
	"artk.dev/asynctx"
	"artk.dev/clone"
	"context"
	"sync"
)

var _ Observer[any] = (&Mux[any]{}).Observe

// Mux is a thread-safe event multiplexer.
type Mux[Event any] struct {
	mutex sync.RWMutex   // 24 bytes on 64 bits.
	wg    sync.WaitGroup // 12 bytes on 64 bits.

	// On 64-bit systems, 4 bytes of padding will be inserted here to
	// ensure that 64-bit words remain aligned.

	observers          []Observer[Event]           // 24 bytes on 64 bits.
	observerMiddleware []ObserverMiddleware[Event] // 24 bytes on 64 bits.
	contextMiddleware  []ContextMiddleware         // 24 bytes on 64 bits.
}

func (m *Mux[Event]) Observe(originalCtx context.Context, event Event) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	l := len(m.observers)
	m.wg.Add(l)

	// We force the creation of a derived context for safety reasons.
	// Since we know that this context will not be cancellable, we can
	// safely share it across all observers.
	asyncSafeContext := asynctx.From(originalCtx)

	// Observers are notified concurrently.
	for i := range l {
		// Apply the context middleware.
		//
		// Note that we know nothing about the particulars of each
		// context middleware being applied. Therefore, we need to
		// provide a different context to each observer to be safe.
		ctx := asyncSafeContext
		for _, middleware := range m.contextMiddleware {
			ctx = middleware(ctx)
		}

		// Apply the observer middleware.
		observer := m.observers[i]
		for _, middleware := range m.observerMiddleware {
			observer = middleware(observer)
		}

		// Prevent shallow copies at the expense of performance.
		event = clone.Of(event)

		// When this goroutine will finish is unspecified, which means
		// that we cannot rely on the mutex to make any operations
		// thread-safe. In practice, this means that all middleware has
		// to have been applied by this point. This precludes
		// optimizations such as applying the middleware inside the
		// goroutine below.
		//
		// Note that using the sync.WaitGroup is still safe even if
		// we are not holding the mutex.
		go func(
			ctx context.Context,
			observer Observer[Event],
			event Event,
		) {
			defer m.wg.Done()

			// Call the observer.
			//
			// While we ultimately ignore the error here, it was
			// made available to the middleware. This can be used,
			// e.g., for logging.
			_ = observer(ctx, event)
		}(ctx, observer, event)
	}

	return nil
}

func (m *Mux[Event]) WillNotify(
	observers ...Observer[Event],
) *Mux[Event] {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.observers = append(m.observers, observers...)

	// Chaining improves DX.
	return m
}

// WithContextMiddleware registers context middleware.
//
// Context middleware will always be applied before observer middleware.
func (m *Mux[Event]) WithContextMiddleware(
	middleware ...ContextMiddleware,
) *Mux[Event] {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.contextMiddleware = append(m.contextMiddleware, middleware...)

	// Chaining improves DX.
	return m
}

// WithObserverMiddleware registers observer middleware.
//
// Context middleware will always be applied before observer middleware.
func (m *Mux[Event]) WithObserverMiddleware(
	middleware ...ObserverMiddleware[Event],
) *Mux[Event] {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.observerMiddleware = append(m.observerMiddleware, middleware...)

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

func New[Event any]() *Mux[Event] {
	return &Mux[Event]{}
}
