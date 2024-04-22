package apperror

import (
	"errors"
	"fmt"
)

type notFoundError struct {
	error
}

func (e notFoundError) NotFound() bool {
	return true
}

func (e notFoundError) Kind() Kind {
	return NotFoundError
}

// NotFound creates a new not found error.
func NotFound(msg string, a ...any) error {
	return &notFoundError{error: fmt.Errorf(msg, a...)}
}

// AsNotFound wraps an existing error as a not found error.
// It returns nil for nil errors.
func AsNotFound(err error) error {
	if err == nil {
		return nil
	}

	return &notFoundError{error: err}
}

// IsNotFound matches not found errors.
func IsNotFound(err error) bool {
	var target interface {
		NotFound() bool
	}
	return errors.As(err, &target) && target.NotFound()
}
