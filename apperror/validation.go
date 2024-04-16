package apperror

import (
	"errors"
	"fmt"
)

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
func Validation(msg string, a ...any) error {
	return &validationErr{error: fmt.Errorf(msg, a...)}
}

// IsValidation checks fit he error is a validation error.
func IsValidation(err error) bool {
	var target interface {
		Validation() bool
	}
	return errors.As(err, &target) && target.Validation()
}
