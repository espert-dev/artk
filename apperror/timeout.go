package apperror

import "errors"

type timeoutErr struct {
	error
}

func (e timeoutErr) Timeout() bool {
	return true
}

// Timeout creates a new timeout error.
func Timeout(msg string) error {
	return &timeoutErr{error: errors.New(msg)}
}

// IsTimeout checks if the error is a timeout error.
func IsTimeout(err error) bool {
	var target interface {
		Timeout() bool
	}
	return errors.As(err, &target) && target.Timeout()
}
