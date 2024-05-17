package eventmux

import "context"

// Observer is a function that can process a specific event type.
//
// The correct handling of errors depends on what the observer is:
//
//   - Event brokers, such as Mux, must indicate whether the propagation of
//     the event succeeded, regardless of whether the observers succeeded.
//
//   - Middleware must propagate the errors they receive unless the
//     purpose of that middleware is to modify the error.
//
//   - Asynchronous operations must return whether they succeeded.
//     Brokers may depend on this to retry failed operations.
type Observer[Event any] func(ctx context.Context, e Event) error
