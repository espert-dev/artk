package httperror

import (
	"artk.dev/apperror"
	"net/http"
)

// EncodeKind maps an apperror.Kind to an HTTP status code.
func EncodeKind(kind apperror.Kind) int {
	switch kind {
	case apperror.OK:
		return http.StatusNoContent
	case apperror.ValidationError:
		return http.StatusBadRequest
	case apperror.UnauthorizedError:
		return http.StatusUnauthorized
	case apperror.ForbiddenError:
		return http.StatusForbidden
	case apperror.NotFoundError:
		return http.StatusNotFound
	case apperror.ConflictError:
		return http.StatusConflict
	case apperror.PreconditionFailedError:
		return http.StatusPreconditionFailed
	case apperror.TooManyRequestsError:
		return http.StatusTooManyRequests
	case apperror.TimeoutError:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// DecodeKind maps an HTTP status code to an apperror.Kind.
//
//gocyclo:ignore
func DecodeKind(status int) apperror.Kind {
	// While it might be advisable to handle success status codes before
	// calling this function, we provide a reasonably safe default.
	if status < http.StatusBadRequest {
		return apperror.OK
	}

	switch status {
	case http.StatusBadRequest:
		return apperror.ValidationError
	case http.StatusUnauthorized:
		return apperror.UnauthorizedError
	case http.StatusForbidden:
		return apperror.ForbiddenError
	case http.StatusNotFound:
		return apperror.NotFoundError
	case http.StatusConflict:
		return apperror.ConflictError
	case http.StatusPreconditionFailed:
		return apperror.PreconditionFailedError
	case http.StatusTooManyRequests:
		return apperror.TooManyRequestsError
	case http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		// Detect infrastructure (e.g., load balancer) errors,
		// in addition to application errors.
		return apperror.TimeoutError
	default:
		// Errors in the 400 range are client errors.
		// Validation is the closest kind.
		if status < http.StatusInternalServerError {
			return apperror.ValidationError
		}

		// Errors in the 500 range are server errors.
		// Unknown is the closest kind.
		return apperror.UnknownError
	}
}
