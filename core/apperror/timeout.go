package apperror

import (
	"errors"
	"fmt"
)

type timeoutErr struct {
	error
}

func (e timeoutErr) Timeout() bool {
	return true
}

func (e timeoutErr) Kind() Kind {
	return TimeoutKind
}

// Timeout creates a new timeout error.
func Timeout(msg string, a ...any) error {
	return &timeoutErr{error: fmt.Errorf(msg, a...)}
}

// IsTimeout checks if the error is a timeout error.
func IsTimeout(err error) bool {
	var target interface {
		Timeout() bool
	}
	return errors.As(err, &target) && target.Timeout()
}
