package testbarrier_test

import (
	"artk.dev/testbarrier"
	"sync"
	"testing"
	"time"
)

func TestWaitForGroup_succeeds_if_done_before_timeout_expires(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	wg.Wait() // Finishes immediately because the counter is zero.
	testbarrier.WaitForGroup(t, &wg, 100*365*24*time.Hour)
}

func TestWaitForGroup_calls_FailNow_if_timeout_expires(t *testing.T) {
	t.Parallel()

	success := make(chan struct{})
	go func() {
		fakeT := &testingT{
			onHelper:  make(chan struct{}),
			onError:   make(chan struct{}),
			onFailNow: make(chan struct{}),
		}

		go func() {
			var wg sync.WaitGroup
			testbarrier.WaitForGroup(fakeT, &wg, time.Nanosecond)
		}()

		<-fakeT.onHelper
		<-fakeT.onError
		<-fakeT.onFailNow
		success <- struct{}{}
	}()

	select {
	case <-success:
		// Hurrah!
	case <-time.NewTicker(5 * time.Second).C:
		t.Errorf("property was not satisfied within timeout")
	}
}
