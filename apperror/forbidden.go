package apperror

import "errors"

type forbiddenErr struct {
	error
}

func (e forbiddenErr) Forbidden() bool {
	return true
}

func (e forbiddenErr) Kind() Kind {
	return ForbiddenKind
}

// Forbidden creates a new forbidden error.
func Forbidden(msg string) error {
	return &forbiddenErr{error: errors.New(msg)}
}

// IsForbidden checks if the error is a forbidden error.
func IsForbidden(err error) bool {
	var target interface {
		Forbidden() bool
	}
	return errors.As(err, &target) && target.Forbidden()
}
