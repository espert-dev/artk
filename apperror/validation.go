package apperror

import "errors"

type validationErr struct {
	error
}

func (e validationErr) Validation() bool {
	return true
}

func (e validationErr) Kind() Kind {
	return ValidationKind
}

// Validation creates a new validation error.
func Validation(msg string) error {
	return &validationErr{error: errors.New(msg)}
}

// IsValidation checks fit he error is a validation error.
func IsValidation(err error) bool {
	var target interface {
		Validation() bool
	}
	return errors.As(err, &target) && target.Validation()
}
