//go:generate go run golang.org/x/tools/cmd/stringer@latest -type Kind .
package apperror

import (
	"errors"
)

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
	// Fast detection, where available.
	var kinder interface {
		Kind() Kind
	}
	if errors.As(err, &kinder) {
		return kinder.Kind()
	}

	// The Go standard library uses Timeout, so put that one first.
	if IsTimeout(err) {
		return TimeoutKind
	}

	// Slow detection. Handle the case where the errors have not been
	// defined with this library. There is a fair chance that by the point
	// we get here we do not have a semantic error at all, in which case
	// the processing speed might matter a bit less.
	if IsNotModified(err) {
		return NotModifiedKind
	} else if IsValidation(err) {
		return ValidationKind
	} else if IsUnauthorized(err) {
		return UnauthorizedKind
	} else if IsForbidden(err) {
		return ForbiddenKind
	} else if IsNotFound(err) {
		return NotFoundKind
	} else if IsConflict(err) {
		return ConflictKind
	} else if IsPreconditionFailed(err) {
		return PreconditionFailedKind
	} else if IsTooManyRequests(err) {
		return TooManyRequestsKind
	}

	return UnknownKind
}
