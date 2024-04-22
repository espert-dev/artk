package apperror

import (
	"errors"
	"fmt"
)

type tooManyRequestsError struct {
	error
}

func (e tooManyRequestsError) TooManyRequests() bool {
	return true
}

func (e tooManyRequestsError) Kind() Kind {
	return TooManyRequestsError
}

// TooManyRequests creates a new too many requests error.
func TooManyRequests(msg string, a ...any) error {
	return tooManyRequestsError{error: fmt.Errorf(msg, a...)}
}

// AsTooManyRequests wraps an existing error as a too many requests error.
// It returns nil for nil errors.
func AsTooManyRequests(err error) error {
	if err == nil {
		return nil
	}

	return tooManyRequestsError{error: err}
}

// IsTooManyRequests matches too many requests errors.
func IsTooManyRequests(err error) bool {
	var target interface {
		TooManyRequests() bool
	}
	return errors.As(err, &target) && target.TooManyRequests()
}
