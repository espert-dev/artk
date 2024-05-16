// Package asynctx provides a mechanism to create derived contexts that can
// be used in an asynchronous setting.
//
// These contexts will share the values of the original context, but not
// Deadline, Error or Done.
package asynctx

import (
	"context"
	"time"
)

var _ context.Context = &from{}

type from struct {
	Parent context.Context
}

func (c from) Deadline() (deadline time.Time, ok bool) {
	// Identical to the implementation of context.Background().
	return
}

func (c from) Done() <-chan struct{} {
	// Identical to the implementation of context.Background().
	return nil
}

func (c from) Err() error {
	// Identical to the implementation of context.Background().
	return nil
}

func (c from) Value(key any) any {
	return c.Parent.Value(key)
}

// From returns a new derived context that can be safely used in asynchronous
// operations that have a different lifetime than the original.
//
// Functionally, it preserves the behavior of the parent context's Value,
// but its Deadline, Done, and Error methods behave like those of
// context.Background.
func From(parent context.Context) context.Context {
	return from{Parent: parent}
}
