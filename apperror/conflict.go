package apperror

import (
	"errors"
	"fmt"
)

type conflictErr struct {
	error
}

func (e conflictErr) Conflict() bool {
	return true
}

func (e conflictErr) Kind() Kind {
	return ConflictKind
}

// Conflict creates a new conflict error.
func Conflict(msg string, a ...any) error {
	return &conflictErr{error: fmt.Errorf(msg, a...)}
}

// IsConflict checks if the error is a conflict error.
func IsConflict(err error) bool {
	var target interface {
		Conflict() bool
	}
	return errors.As(err, &target) && target.Conflict()
}
