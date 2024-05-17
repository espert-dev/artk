package testbarrier

type testingT interface {
	Error(args ...any)
	FailNow()
	Helper()
}
