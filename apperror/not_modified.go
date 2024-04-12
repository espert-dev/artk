package apperror

import "errors"

type notModifiedErr struct {
	error
}

func (e notModifiedErr) NotModified() bool {
	return true
}

func (e notModifiedErr) Kind() Kind {
	return NotModifiedKind
}

// NotModified creates a new not modified error.
// This is usually a sentinel error and does not indicate a problem.
func NotModified(msg string) error {
	return &notModifiedErr{error: errors.New(msg)}
}

// IsNotModified checks if the error is a not modified error.
func IsNotModified(err error) bool {
	var target interface {
		NotModified() bool
	}
	return errors.As(err, &target) && target.NotModified()
}
