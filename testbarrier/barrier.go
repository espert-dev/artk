// Package testbarrier implements a barrier for testing asynchronous systems.
package testbarrier

import (
	"artk.dev/assume"
	"sync"
	"time"
)

// Barrier blocks tests for a limited duration until an event happens.
type Barrier struct {
	ch        chan struct{}
	closeOnce sync.Once
}

// Lift the barrier. Must be called when the event happens.
//
// This method is idempotent and can be called multiple times.
func (b *Barrier) Lift() {
	b.closeOnce.Do(func() {
		close(b.ch)
	})
}

// Wait for the barrier to lift indefinitely.
//
// Eventually, the go test timeout will kick in.
// It can introduce higher delays than WaitFor on failures, but on the other
// hand it is much friendlier to debugging.
func (b *Barrier) Wait() {
	<-b.ch
}

// WaitFor for the barrier to lift for up to a duration `d`.
// If the deadline expires, the test will fail immediately.
func (b *Barrier) WaitFor(t testingT, d time.Duration) {
	assume.NotZero(t)

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
func New() *Barrier {
	return &Barrier{ch: make(chan struct{})}
}
