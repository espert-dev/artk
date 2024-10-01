package event

import (
	"artk.dev/assume"
	"artk.dev/asynctx"
	"artk.dev/clone"
	"artk.dev/ptr"
	"context"
	"sync"
)

var _ Observer[any] = (&Stream[any]{}).Observe

// Stream is a thread-safe in-memory event stream.
type Stream[Event any] struct {
	consumerChannels  []chan<- eventMsg[Event] // 24 bytes on 64 bits.
	mutex             sync.RWMutex             // 24 bytes on 64 bits.
	consumerWaitGroup sync.WaitGroup           // 12 bytes on 64 bits.

	// Equals (queueSize - defaultQueueSize).
	// This way, the zero value of Stream will use the default value.
	extraQueueSize int32 //  4 bytes.

	// Up to here, 64 bytes (common cache line size) on 64 bits.

	contextMiddleware  []ContextMiddleware         // 24 bytes on 64 bits.
	observerMiddleware []ObserverMiddleware[Event] // 24 bytes on 64 bits.
}

// Observe an event and propagate it to existing consumers.
//
// This function never returns an error.
func (s *Stream[Event]) Observe(ctx context.Context, e Event) error {
	if ctx.Err() != nil {
		// The context was cancelled: do not call observers.
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, output := range s.consumerChannels {
		// Trade performance for safety: prevent shallow copies.
		msg := eventMsg[Event]{
			ctx:   asynctx.From(ctx),
			event: clone.Of(e),
		}

		select {
		case output <- msg:
			// Message sent successfully.
		default:
			// Queue full: message dropped.
		}
	}

	return nil
}

// WillNotify and handle events.
//
// Events produced before this function is called cannot be observed.
// The consume function runs in a new goroutine.
// This function never returns an error.
func (s *Stream[Event]) WillNotify(consume Observer[Event]) *Stream[Event] {
	// Support shutdown.
	s.consumerWaitGroup.Add(1)
	ch := s.newConsumerChannel()

	go func() {
		defer s.consumerWaitGroup.Done()

		// Consume messages.
		for msg := range ch {
			ctx := msg.ctx

			// Apply context middleware.
			for _, middleware := range s.contextMiddleware {
				ctx = middleware(ctx)
			}

			// Apply the observer middleware.
			h := consume
			for _, middleware := range s.observerMiddleware {
				h = middleware(h)
			}

			// Ignore errors for now.
			_ = h(ctx, msg.event)
		}
	}()

	// Chaining improves DX.
	return s
}

// WithContextMiddleware registers context middleware.
//
// Context middleware will always be applied before observer middleware.
func (s *Stream[Event]) WithContextMiddleware(
	middleware ...ContextMiddleware,
) *Stream[Event] {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.contextMiddleware = append(s.contextMiddleware, middleware...)

	// Chaining improves DX.
	return s
}

// WithObserverMiddleware registers observer middleware.
//
// Context middleware will always be applied before observer middleware.
func (s *Stream[Event]) WithObserverMiddleware(
	middleware ...ObserverMiddleware[Event],
) *Stream[Event] {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.observerMiddleware = append(s.observerMiddleware, middleware...)

	// Chaining improves DX.
	return s
}

// Shutdown the Stream and communicate finishing via the sync.WaitGroup.
func (s *Stream[Event]) Shutdown(wg *sync.WaitGroup) {
	// Synchronously prevent new messages from being sent.
	s.stopEventPropagation()

	// Asynchronously wait for tasks to finish.
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.consumerWaitGroup.Wait()
	}()
}

func (s *Stream[Event]) newConsumerChannel() chan eventMsg[Event] {
	// Memory allocation doesn't need the lock.
	ch := make(chan eventMsg[Event], int(defaultQueueSize+s.extraQueueSize))

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.consumerChannels = append(s.consumerChannels, ch)
	return ch
}

func (s *Stream[Event]) stopEventPropagation() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Closing the channels to signal termination to the ConsumerGroup's.
	for _, output := range s.consumerChannels {
		close(output)
	}

	// Prevent production of new messages.
	s.consumerChannels = nil
}

// NewStream creates a Stream with the specified maximum queue size.
func NewStream[Event any](
	optionsFn ...func(options *streamOptions),
) *Stream[Event] {
	options := &streamOptions{
		queueSize: ptr.To(defaultQueueSize),
	}
	for _, fn := range optionsFn {
		fn(options)
	}

	assume.Truef(
		*options.queueSize >= 0,
		"queue size cannot be negative (was %v)",
		*options.queueSize,
	)

	return &Stream[Event]{
		extraQueueSize: *options.queueSize - defaultQueueSize,
	}
}

func WithStreamQueueSize(size int32) func(options *streamOptions) {
	return func(options *streamOptions) {
		options.queueSize = &size
	}
}

type streamOptions struct {
	queueSize *int32
}

type eventMsg[Event any] struct {
	ctx   context.Context
	event Event
}

const defaultQueueSize int32 = 128
