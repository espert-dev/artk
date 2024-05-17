package testbarrier_test

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
