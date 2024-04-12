package apperror

import "errors"

type unauthorizedErr struct {
	error
}

func (e unauthorizedErr) Unauthorized() bool {
	return true
}

// Unauthorized creates a new unauthorized error.
func Unauthorized(msg string) error {
	return &unauthorizedErr{error: errors.New(msg)}
}

// IsUnauthorized checks if the error is a unauthorized error.
func IsUnauthorized(err error) bool {
	var target interface {
		Unauthorized() bool
	}
	return errors.As(err, &target) && target.Unauthorized()
}
