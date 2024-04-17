package testbarrier_test

import (
	"artk.dev/core/testbarrier"
	"testing"
	"time"
)

func TestBarrier_Lift(t *testing.T) {
	barrier := testbarrier.New()
	go barrier.Lift()
	barrier.Wait(t, 100*365*24*time.Hour)
}

func TestBarrier_Wait_expires(t *testing.T) {
	success := make(chan struct{})
	go func() {
		fakeT := &testingT{
			onHelper:  make(chan struct{}),
			onError:   make(chan struct{}),
			onFailNow: make(chan struct{}),
		}

		go func() {
			barrier := testbarrier.New()
			barrier.Wait(fakeT, time.Nanosecond)
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

type testingT struct {
	onHelper  chan struct{}
	onError   chan struct{}
	onFailNow chan struct{}
}

func (t *testingT) Helper() {
	t.onHelper <- struct{}{}
}

func (t *testingT) Error(_ ...any) {
	t.onError <- struct{}{}
}

func (t *testingT) FailNow() {
	t.onFailNow <- struct{}{}
}
