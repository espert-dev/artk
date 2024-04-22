// Package testbarrier implements a barrier for testing asynchronous systems.
package testbarrier

import (
	"artk.dev/assume"
	"time"
)

// Barrier blocks tests for a limited duration until an event happens.
type Barrier struct {
	ch chan struct{}
}

// Lift the barrier. Must be called when the event happens.
func (b Barrier) Lift() {
	b.ch <- struct{}{}
}

// Wait for the barrier to lift for up to a duration `d`.
// If the deadline expires, the test will fail immediately.
func (b Barrier) Wait(t testingT, d time.Duration) {
	assume.NotNil(t)

	t.Helper()

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	select {
	case <-b.ch:
		// All good.
	case <-ticker.C:
		t.Error("barrier timeout exceeded")
		t.FailNow()
	}
}

// New creates a Barrier.
func New() Barrier {
	return Barrier{ch: make(chan struct{})}
}

type testingT interface {
	Error(args ...any)
	FailNow()
	Helper()
}
