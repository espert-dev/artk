package apperror

import (
	"errors"
	"fmt"
)

type tooManyRequestsErr struct {
	error
}

func (e tooManyRequestsErr) TooManyRequests() bool {
	return true
}

func (e tooManyRequestsErr) Kind() Kind {
	return TooManyRequestsKind
}

// TooManyRequests creates a new too many requests error.
func TooManyRequests(msg string, a ...any) error {
	return &tooManyRequestsErr{error: fmt.Errorf(msg, a...)}
}

// IsTooManyRequests checks if the error is a too many requests error.
func IsTooManyRequests(err error) bool {
	var target interface {
		TooManyRequests() bool
	}
	return errors.As(err, &target) && target.TooManyRequests()
}
