package apperror

import (
	"fmt"
)

// As wraps an existing error into a Kind.
// Invalid Kind values are mapped to Unknown.
// If the kind is OK or the error is nil, the function will return nil.
//
//gocyclo:ignore
func As(kind Kind, err error) error {
	if err == nil {
		return nil
	}

	switch kind {
	case OK:
		return nil
	case ValidationError:
		return validationError{error: err}
	case UnauthorizedError:
		return unauthorizedError{error: err}
	case ForbiddenError:
		return forbiddenError{error: err}
	case NotFoundError:
		return notFoundError{error: err}
	case ConflictError:
		return conflictError{error: err}
	case PreconditionFailedError:
		return preconditionFailedError{error: err}
	case TooManyRequestsError:
		return tooManyRequestsError{error: err}
	case TimeoutError:
		return timeoutError{error: err}
	default:
		return unknownError{error: err}
	}
}

// New creates an error of the specified Kind.
// Invalid Kind values are mapped to Unknown.
// If the kind is OK, the function will return nil.
func New(kind Kind, msg string, a ...any) error {
	err := fmt.Errorf(msg, a...)
	return As(kind, err)
}

// IsFinal returns true for error kinds that are not expected to change by
// merely repeating the operation.
//
// Note that this function exclusively considers error kinds, not the context
// in which they were produced. Not all operations can be repeated safely, and
// this function cannot provide guidance on making such a judgment.
func IsFinal(err error) bool {
	switch kind := KindOf(err); kind {
	case UnknownError,
		TooManyRequestsError,
		TimeoutError:
		return false
	default:
		return true
	}
}
