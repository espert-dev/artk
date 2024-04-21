package apperror

import (
	"errors"
	"fmt"
)

type timeoutError struct {
	error
}

func (e timeoutError) Timeout() bool {
	return true
}

func (e timeoutError) Kind() Kind {
	return TimeoutError
}

// Timeout creates a new timeout error.
func Timeout(msg string, a ...any) error {
	return &timeoutError{error: fmt.Errorf(msg, a...)}
}

// AsTimeout wraps an existing error as a timeout error.
// It returns nil for nil errors.
func AsTimeout(err error) error {
	if err == nil {
		return nil
	}

	return &timeoutError{error: err}
}

// IsTimeout matches timeout errors.
func IsTimeout(err error) bool {
	var target interface {
		Timeout() bool
	}
	return errors.As(err, &target) && target.Timeout()
}
