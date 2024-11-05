package testbarrier_test

import (
	"artk.dev/testbarrier"
	"testing"
	"time"
)

func TestBarrier_Wait(t *testing.T) {
	t.Parallel()

	barrier := testbarrier.New()
	go barrier.Lift()
	barrier.Wait(t)
}

func TestBarrier_Wait_can_wait_multiple_times(t *testing.T) {
	t.Parallel()

	barrier := testbarrier.New()
	go barrier.Lift()

	barrier.Wait(t)
	barrier.Wait(t)
	barrier.Wait(t)
}

func TestBarrier_WaitFor_can_wait_multiple_times(t *testing.T) {
	t.Parallel()

	barrier := testbarrier.New()
	go barrier.Lift()

	barrier.WaitFor(t, 100*365*24*time.Hour)
	barrier.WaitFor(t, 100*365*24*time.Hour)
	barrier.WaitFor(t, 100*365*24*time.Hour)
}

func TestBarrier_WaitFor_can_lift_multiple_times(t *testing.T) {
	t.Parallel()

	barrier := testbarrier.New()
	barrier.Lift()
	barrier.Lift()
	barrier.Lift()

	barrier.WaitFor(t, 100*365*24*time.Hour)
}

func TestBarrier_WaitFor_can_lift_and_wait_multiple_times(t *testing.T) {
	t.Parallel()

	barrier := testbarrier.New()
	barrier.Lift()
	barrier.Lift()
	barrier.Lift()

	barrier.WaitFor(t, 100*365*24*time.Hour)
	barrier.WaitFor(t, 100*365*24*time.Hour)
	barrier.WaitFor(t, 100*365*24*time.Hour)
}

func TestBarrier_WaitFor_ok_if_lifted_before_timeout_expires(t *testing.T) {
	t.Parallel()

	barrier := testbarrier.New()
	go barrier.Lift()
	barrier.WaitFor(t, 100*365*24*time.Hour)
}

func TestBarrier_WaitFor_calls_FailNow_if_timeout_expires(t *testing.T) {
	t.Parallel()

	success := make(chan struct{})
	go func() {
		fakeT := &testingT{
			onHelper:  make(chan struct{}),
			onError:   make(chan struct{}),
			onFailNow: make(chan struct{}),
		}

		go func() {
			// Nothing can lift the barrier.
			// This guarantees that it will time out.
			barrier := testbarrier.New()
			barrier.WaitFor(fakeT, time.Nanosecond)
		}()

		<-fakeT.onHelper
		<-fakeT.onError
		<-fakeT.onFailNow
		success <- struct{}{}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	select {
	case <-success:
		// Hurrah!
	case <-ticker.C:
		t.Errorf("property was not satisfied within timeout")
	}
}

func TestBarrier_WaitFor_never_blocks_after_Lift(t *testing.T) {
	t.Parallel()

	barrier := testbarrier.New()
	go barrier.Lift()

	for range 100 {
		barrier.WaitFor(t, 100*365*24*time.Hour)
	}
}
