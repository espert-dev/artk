package testbarrier

import (
	"artk.dev/assume"
	"sync"
	"time"
)

// WaitForGroup waits for a sync.WaitGroup for up to the specified duration.
// If the deadline expires, the test will fail immediately.
func WaitForGroup(t testingT, wg *sync.WaitGroup, d time.Duration) {
	assume.NotZero(t)
	assume.NotZero(wg)

	t.Helper()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		t.Error("group timeout exceeded")
		t.FailNow()
	case <-done:
		// Success!
	}
}
