package apperror

import "errors"

type conflictErr struct {
	error
}

func (e conflictErr) Conflict() bool {
	return true
}

// Conflict creates a new conflict error.
func Conflict(msg string) error {
	return &conflictErr{error: errors.New(msg)}
}

// IsConflict checks if the error is a conflict error.
func IsConflict(err error) bool {
	var target interface {
		Conflict() bool
	}
	return errors.As(err, &target) && target.Conflict()
}
