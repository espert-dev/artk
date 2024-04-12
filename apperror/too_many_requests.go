package apperror

import "errors"

type tooManyRequestsErr struct {
	error
}

func (e tooManyRequestsErr) TooManyRequests() bool {
	return true
}

// TooManyRequests creates a new tooManyRequests error.
func TooManyRequests(msg string) error {
	return &tooManyRequestsErr{error: errors.New(msg)}
}

// IsTooManyRequests checks if the error is a tooManyRequests error.
func IsTooManyRequests(err error) bool {
	var target interface {
		TooManyRequests() bool
	}
	return errors.As(err, &target) && target.TooManyRequests()
}
