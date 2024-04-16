package apperror

import (
	"errors"
	"fmt"
)

type unauthorizedErr struct {
	error
}

func (e unauthorizedErr) Unauthorized() bool {
	return true
}

func (e unauthorizedErr) Kind() Kind {
	return UnauthorizedKind
}

// Unauthorized creates a new unauthorized error.
func Unauthorized(msg string, a ...any) error {
	return &unauthorizedErr{error: fmt.Errorf(msg, a...)}
}

// IsUnauthorized checks if the error is a unauthorized error.
func IsUnauthorized(err error) bool {
	var target interface {
		Unauthorized() bool
	}
	return errors.As(err, &target) && target.Unauthorized()
}
