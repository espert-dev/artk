package require

import (
	"artk.dev/internal/assert"
)

type TestingT interface {
	Helper()
	Errorf(format string, args ...any)
	FailNow()
}

func As[V any](t TestingT, x any) V {
	v, ok := assert.As[V](t, x)
	if !ok {
		t.FailNow()
	}

	return v
}

func Nil(t TestingT, got any) {
	t.Helper()

	if !assert.Nil(t, got) {
		t.FailNow()
	}
}

func NotNil(t TestingT, got any) {
	t.Helper()

	if !assert.NotNil(t, got) {
		t.FailNow()
	}
}

func Equal[V comparable](t TestingT, expected, got V) {
	t.Helper()

	if !assert.Equal(t, expected, got) {
		t.FailNow()
	}
}

func NotEqual[V comparable](t TestingT, x, y V) {
	t.Helper()

	if !assert.NotEqual(t, x, y) {
		t.FailNow()
	}
}

func Same[V any](t TestingT, expected, got V) {
	t.Helper()

	if !assert.Same(t, expected, got) {
		t.FailNow()
	}
}

func NotSame[V any](t TestingT, x, y V) {
	t.Helper()

	if !assert.NotSame(t, x, y) {
		t.FailNow()
	}
}
