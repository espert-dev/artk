package apperror

import (
	"errors"
	"fmt"
)

type conflictError struct {
	error
}

func (e conflictError) Conflict() bool {
	return true
}

func (e conflictError) Kind() Kind {
	return ConflictError
}

// Conflict creates a new conflict error.
func Conflict(msg string, a ...any) error {
	return &conflictError{error: fmt.Errorf(msg, a...)}
}

// AsConflict wraps an existing error as a conflict error.
// It returns nil for nil errors.
func AsConflict(err error) error {
	if err == nil {
		return nil
	}

	return &conflictError{error: err}
}

// IsConflict matches conflict errors.
func IsConflict(err error) bool {
	var target interface {
		Conflict() bool
	}
	return errors.As(err, &target) && target.Conflict()
}
