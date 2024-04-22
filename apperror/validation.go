package apperror

import (
	"errors"
	"fmt"
)

type validationError struct {
	error
}

func (e validationError) Validation() bool {
	return true
}

func (e validationError) Kind() Kind {
	return ValidationError
}

// Validation creates a new validation error.
func Validation(msg string, a ...any) error {
	return validationError{error: fmt.Errorf(msg, a...)}
}

// AsValidation wraps an existing error as a validation error.
// It returns nil for nil errors.
func AsValidation(err error) error {
	if err == nil {
		return nil
	}

	return validationError{error: err}
}

// IsValidation matches validation errors.
func IsValidation(err error) bool {
	var target interface {
		Validation() bool
	}
	return errors.As(err, &target) && target.Validation()
}
