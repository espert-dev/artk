// Package grpcerror provides serialization and deserialization of errors
// that follow apperror conventions.
//
// When in doubt about code mappings, check this:
// https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md
package grpcerror

import (
	"artk.dev/apperror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Encode an application error into a gRPC error.
func Encode(err error) error {
	if err == nil {
		return nil
	}

	kind := apperror.KindOf(err)
	code := EncodeKind(kind)
	return status.Error(code, err.Error())
}

// Decode a gRPC error into an application error.
func Decode(err error) error {
	if err == nil {
		return nil
	}

	s, ok := status.FromError(err)
	if !ok {
		// Will be handled as an unknown error.
		return apperror.Unknown("cannot parse gRPC error: %w", err)
	}

	kind := DecodeKind(s.Code())
	return apperror.New(kind, s.Message())
}

// EncodeKind encodes an apperror.Kind into a gRPC codes.Code.
func EncodeKind(kind apperror.Kind) codes.Code {
	switch kind {
	case apperror.OK:
		return codes.OK
	case apperror.ValidationError:
		return codes.InvalidArgument
	case apperror.UnauthorizedError:
		return codes.Unauthenticated
	case apperror.ForbiddenError:
		return codes.PermissionDenied
	case apperror.NotFoundError:
		return codes.NotFound
	case apperror.ConflictError:
		return codes.AlreadyExists
	case apperror.PreconditionFailedError:
		return codes.FailedPrecondition
	case apperror.TooManyRequestsError:
		return codes.Unavailable
	case apperror.TimeoutError:
		return codes.DeadlineExceeded
	default:
		return codes.Unknown
	}
}

// DecodeKind decodes a codes.Code into an apperror.Kind.
func DecodeKind(code codes.Code) apperror.Kind {
	switch code {
	case codes.OK:
		return apperror.OK
	case codes.InvalidArgument:
		return apperror.ValidationError
	case codes.Unauthenticated:
		return apperror.UnauthorizedError
	case codes.PermissionDenied:
		return apperror.ForbiddenError
	case codes.NotFound:
		return apperror.NotFoundError
	case codes.AlreadyExists:
		return apperror.ConflictError
	case codes.FailedPrecondition:
		return apperror.PreconditionFailedError
	case codes.Unavailable:
		return apperror.TooManyRequestsError
	case codes.DeadlineExceeded:
		return apperror.TimeoutError
	default:
		return apperror.UnknownError
	}
}
