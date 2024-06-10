package apperror

import (
	"errors"
	"fmt"
)

type forbiddenError struct {
	error
}

func (e forbiddenError) Forbidden() bool {
	return true
}

func (e forbiddenError) Kind() Kind {
	return ForbiddenError
}

// Forbidden creates a new forbidden error.
func Forbidden(msg string) error {
	return forbiddenError{error: errors.New(msg)}
}

// Forbiddenf creates a new forbidden error.
func Forbiddenf(msg string, a ...any) error {
	return forbiddenError{error: fmt.Errorf(msg, a...)}
}

// AsForbidden wraps an existing error as a forbidden error.
// It returns nil for nil errors.
func AsForbidden(err error) error {
	if err == nil {
		return nil
	}

	return forbiddenError{error: err}
}

// IsForbidden matches forbidden errors.
func IsForbidden(err error) bool {
	var target interface {
		Forbidden() bool
	}
	return errors.As(err, &target) && target.Forbidden()
}
