package apperror

import "errors"

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
func TooManyRequests(msg string) error {
	return &tooManyRequestsErr{error: errors.New(msg)}
}

// IsTooManyRequests checks if the error is a too many requests error.
func IsTooManyRequests(err error) bool {
	var target interface {
		TooManyRequests() bool
	}
	return errors.As(err, &target) && target.TooManyRequests()
}
