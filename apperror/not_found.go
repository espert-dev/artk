package apperror

import (
	"errors"
	"fmt"
)

type notFoundErr struct {
	error
}

func (e notFoundErr) NotFound() bool {
	return true
}

func (e notFoundErr) Kind() Kind {
	return NotFoundKind
}

// NotFound creates a new not found error.
func NotFound(msg string, a ...any) error {
	return &notFoundErr{error: fmt.Errorf(msg, a...)}
}

// IsNotFound checks if the error is a not found error.
func IsNotFound(err error) bool {
	var target interface {
		NotFound() bool
	}
	return errors.As(err, &target) && target.NotFound()
}
