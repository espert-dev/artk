package apperror

import (
	"errors"
	"fmt"
)

type preconditionFailedError struct {
	error
}

func (e preconditionFailedError) PreconditionFailed() bool {
	return true
}

func (e preconditionFailedError) Kind() Kind {
	return PreconditionFailedError
}

// PreconditionFailed creates a new precondition failed error.
func PreconditionFailed(msg string, a ...any) error {
	return preconditionFailedError{error: fmt.Errorf(msg, a...)}
}

// AsPreconditionFailed wraps an existing error as a precondition failed error.
// It returns nil for nil errors.
func AsPreconditionFailed(err error) error {
	if err == nil {
		return nil
	}

	return preconditionFailedError{error: err}
}

// IsPreconditionFailed matches precondition failed errors.
func IsPreconditionFailed(err error) bool {
	var target interface {
		PreconditionFailed() bool
	}
	return errors.As(err, &target) && target.PreconditionFailed()
}
