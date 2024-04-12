package apperror

import "errors"

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
func PreconditionFailed(msg string) error {
	return &preconditionFailedErr{error: errors.New(msg)}
}

// IsPreconditionFailed checks if the error is a precondition failed error.
func IsPreconditionFailed(err error) bool {
	var target interface {
		PreconditionFailed() bool
	}
	return errors.As(err, &target) && target.PreconditionFailed()
}
