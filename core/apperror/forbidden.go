package apperror

import (
	"errors"
	"fmt"
)

type forbiddenErr struct {
	error
}

func (e forbiddenErr) Forbidden() bool {
	return true
}

func (e forbiddenErr) Kind() Kind {
	return ForbiddenKind
}

// Forbidden creates a new forbidden error.
func Forbidden(msg string, a ...any) error {
	return &forbiddenErr{error: fmt.Errorf(msg, a...)}
}

// IsForbidden checks if the error is a forbidden error.
func IsForbidden(err error) bool {
	var target interface {
		Forbidden() bool
	}
	return errors.As(err, &target) && target.Forbidden()
}
