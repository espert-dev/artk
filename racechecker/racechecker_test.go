package racechecker_test

type T struct {
	NumCallsToHelper int
	NumCallsToSkip   int
}

func (t *T) Helper() {
	t.NumCallsToHelper++
}

func (t *T) Skip(_ ...any) {
	t.NumCallsToSkip++
}
