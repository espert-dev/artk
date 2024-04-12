package apperror

import "errors"

type preconditionFailedErr struct {
	error
}

func (e preconditionFailedErr) PreconditionFailed() bool {
	return true
}

// PreconditionFailed creates a new preconditionFailed error.
func PreconditionFailed(msg string) error {
	return &preconditionFailedErr{error: errors.New(msg)}
}

// IsPreconditionFailed checks if the error is a preconditionFailed error.
func IsPreconditionFailed(err error) bool {
	var target interface {
		PreconditionFailed() bool
	}
	return errors.As(err, &target) && target.PreconditionFailed()
}
