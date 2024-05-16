package eventmux

import "context"

// Observer is a function that can process this event.
//
// The interpretation of the error is not trivial and must be done with care.
// The error indicates the result of the operation closest to the caller.
// In particular, it does not indicate the success of successive operations
// that might be triggered by the observer,
// because doing so could lead to unfeasibly long dependency chains,
// or even the impossibility of resolution in the presence of cycles.
// Replacing the error with a promise (func() error), would not help either.
//
// As a partial exception, middleware is generally expected to preserve the
// error of the next middleware or operation in the chain.
type Observer[Event any] func(ctx context.Context, e Event) error

// ObserverMiddleware can be used to run arbitrary actions before and after
// the completion of next.
//
// Note that there are two different ways in which we can use this within a
// Mux:
//
//   - With Mux.WithObserverMiddleware, which will apply it to all observers.
//   - By wrapping the observer before passing it to Mux.WillNotify,
//     which will only apply to this particular observer.
//
// While more involved, the latter approach is more powerful since it allows
// handling based on both the event and the observer.
// For example, logging middleware could attach an observer name or ID to each
// log, making them easier to understand.
type ObserverMiddleware[Event any] func(next Observer[Event]) Observer[Event]

// ContextMiddleware is used to transform the context provided to a function.
//
// It exists as an optimization over ObserverMiddleware whenever the value or
// type of the event is irrelevant, which supports reuse without needlessly
// depending on the instantiation of generics.
type ContextMiddleware func(ctx context.Context) context.Context
