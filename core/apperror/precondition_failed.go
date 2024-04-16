package apperror

import (
	"errors"
	"fmt"
)

type preconditionFailedErr struct {
	error
}

func (e preconditionFailedErr) PreconditionFailed() bool {
	return true
}

func (e preconditionFailedErr) Kind() Kind {
	return PreconditionFailedKind
}

// PreconditionFailed creates a new precondition failed error.
func PreconditionFailed(msg string, a ...any) error {
	return &preconditionFailedErr{error: fmt.Errorf(msg, a...)}
}

// IsPreconditionFailed checks if the error is a precondition failed error.
func IsPreconditionFailed(err error) bool {
	var target interface {
		PreconditionFailed() bool
	}
	return errors.As(err, &target) && target.PreconditionFailed()
}
