package apperror

import (
	"errors"
	"fmt"
)

type unauthorizedError struct {
	error
}

func (e unauthorizedError) Unauthorized() bool {
	return true
}

func (e unauthorizedError) Kind() Kind {
	return UnauthorizedError
}

// Unauthorized creates a new unauthorized error.
func Unauthorized(msg string) error {
	return unauthorizedError{error: errors.New(msg)}
}

// Unauthorizedf creates a new unauthorized error.
func Unauthorizedf(msg string, a ...any) error {
	return unauthorizedError{error: fmt.Errorf(msg, a...)}
}

// AsUnauthorized wraps an existing error as a unauthorized error.
// It returns nil for nil errors.
func AsUnauthorized(err error) error {
	if err == nil {
		return nil
	}

	return unauthorizedError{error: err}
}

// IsUnauthorized matches unauthorized errors.
func IsUnauthorized(err error) bool {
	var target interface {
		Unauthorized() bool
	}
	return errors.As(err, &target) && target.Unauthorized()
}
