//go:generate go run golang.org/x/tools/cmd/stringer@latest -type Kind .

package apperror

import (
	"errors"
)

// Kind implies semantic connotations about an error. In most cases, knowing
// the kind of error is all that user code needs for correct handling,
// and there is no need to consider the exact error type or message.
//
// The numerical values of these constants are not guaranteed to be stable
// and therefore must not be relied on.
type Kind int

const (
	OK Kind = iota
	UnknownError
	ValidationError
	UnauthorizedError
	ForbiddenError
	NotFoundError
	ConflictError
	PreconditionFailedError
	TooManyRequestsError
	TimeoutError
)

// KindValues returns the set of known values of Kind.
func KindValues() []Kind {
	return []Kind{
		OK,
		UnknownError,
		ValidationError,
		UnauthorizedError,
		ForbiddenError,
		NotFoundError,
		ConflictError,
		PreconditionFailedError,
		TooManyRequestsError,
		TimeoutError,
	}
}

// KindOf returns the Kind of an error.
// If the error is nil, it will return OK.
//
// This is usually faster than calling multiple error kind matchers.
func KindOf(err error) Kind {
	// While not essential, supporting OK allows user code to handle
	// success and multiple error kinds with a single switch statement.
	if err == nil {
		return OK
	}

	// Fast detection, where available.
	var kinder interface {
		Kind() Kind
	}
	if errors.As(err, &kinder) {
		return kinder.Kind()
	}

	// The Go standard library uses Timeout, so check it first.
	if IsTimeout(err) {
		return TimeoutError
	}

	return UnknownError
}
