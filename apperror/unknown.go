package apperror

import (
	"errors"
	"fmt"
)

// unknownError is not a semantic error, but it implements the method Kind.
// This allows for faster checks in some places.
type unknownError struct {
	error
}

func (e unknownError) Kind() Kind {
	return UnknownError
}

// Unknown returns a semantic error of UnknownError.
//
// While any non-semantic error will be detected as an unknown error, the
// error returned by this stringConstructor implements the kinder interface and
// can be checked faster.
func Unknown(msg string) error {
	return unknownError{error: errors.New(msg)}
}

// Unknownf returns a semantic error of UnknownError.
//
// While any non-semantic error will be detected as an unknown error, the
// error returned by this stringConstructor implements the kinder interface and
// can be checked faster.
func Unknownf(msg string, a ...any) error {
	return unknownError{error: fmt.Errorf(msg, a...)}
}

// AsUnknown wraps an existing error as a unknown error.
// It returns nil for nil errors.
func AsUnknown(err error) error {
	if err == nil {
		return nil
	}

	return unknownError{error: err}
}

// IsUnknown matches unknown errors.
func IsUnknown(err error) bool {
	if err == nil {
		return false
	}

	kind := KindOf(err)
	return kind == UnknownError
}
