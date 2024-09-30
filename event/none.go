package event

import "context"

// None is an event handler that does nothing.
//
// If you have a service that supports listeners, inject None by default during
// construction, and you can stop checking for `nil` before notifying them.
//
// Example:
//
//	// Construction.
//	h.listener = event.None[MyEvent]
//
//	// Usage (no need to check if `listener` is nil).
//	h.listener(ctx, e)
func None[Event any](_ context.Context, _ Event) error {
	// Deliberately do nothing.
	return nil
}

var _ Observer[any] = None[any]
