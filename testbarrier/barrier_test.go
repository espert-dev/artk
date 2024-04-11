package testbarrier_test

import (
	"github.com/jespert/artk/testbarrier"
	"testing"
	"time"
)

func TestBarrier_Lift(t *testing.T) {
	barrier := testbarrier.New()
	go barrier.Lift()
	barrier.Wait(t, 100*365*24*time.Hour)
}

func TestBarrier_Wait_expires(t *testing.T) {
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
