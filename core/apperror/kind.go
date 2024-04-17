//go:generate go run golang.org/x/tools/cmd/stringer@latest -type Kind .
package apperror

import "errors"

// Kind indicates the category of error.
// There is one for each Is* function in this package.
type Kind int

// The values below don't really matter, but they have been chosen to match
// the HTTP codes because a lot of developers are familiar with them.
//
// The choice for TimeoutKind is a bit tricky, but I want to be able to
// distinguish between a timeout and other kinds of server errors, which would
// be hard if we went with 500.
const (
	UnknownKind            Kind = 0
	NotModifiedKind        Kind = 304
	ValidationKind         Kind = 400
	UnauthorizedKind       Kind = 401
	ForbiddenKind          Kind = 403
	NotFoundKind           Kind = 404
	ConflictKind           Kind = 409
	PreconditionFailedKind Kind = 412
	TooManyRequestsKind    Kind = 429
	TimeoutKind            Kind = 504
)

// KindOf allows potentially faster checking of multiple error types.
func KindOf(err error) Kind {
	kind, ok := fastKindOf(err)
	if ok {
		return kind
	}

	return slowKindOf(err)
}

// fastKindOf quickly computes the Kind of an error created by artk or the
// Go standard library, but it is not compatible with errors created by other
// libraries.
func fastKindOf(err error) (Kind, bool) {
	// Fast detection, where available.
	var kinder interface {
		Kind() Kind
	}
	if errors.As(err, &kinder) {
		return kinder.Kind(), true
	}

	// The Go standard library uses Timeout, so check it first.
	if IsTimeout(err) {
		return TimeoutKind, true
	}

	return UnknownKind, false
}

func slowKindOf(err error) Kind {
	// Check types other than Timeout, which is checked in fastKindOf.
	switch {
	case IsNotModified(err):
		return NotModifiedKind
	case IsValidation(err):
		return ValidationKind
	case IsUnauthorized(err):
		return UnauthorizedKind
	case IsForbidden(err):
		return ForbiddenKind
	case IsNotFound(err):
		return NotFoundKind
	case IsConflict(err):
		return ConflictKind
	case IsPreconditionFailed(err):
		return PreconditionFailedKind
	case IsTooManyRequests(err):
		return TooManyRequestsKind
	default:
		return UnknownKind
	}
}
