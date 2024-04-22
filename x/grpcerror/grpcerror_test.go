package grpcerror_test

import (
	"artk.dev/apperror"
	"artk.dev/x/grpcerror"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestEncode_encodes_kind_into_code(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		t.Run(kind.String(), func(t *testing.T) {
			originalErr := apperror.New(kind, errorMessage)
			encodedErr := grpcerror.Encode(originalErr)

			grpcStatus := asGRPCError(t, encodedErr)
			expected := grpcerror.EncodeKind(kind)
			assertCodeIs(t, grpcStatus, expected)
		})
	}
}

func TestEncode_encodes_message_for_error_kinds(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		if kind == apperror.OK {
			// nil errors do not have messages.
			continue
		}

		t.Run(kind.String(), func(t *testing.T) {
			originalErr := apperror.New(kind, errorMessage)
			encodedErr := grpcerror.Encode(originalErr)

			grpcStatus := asGRPCError(t, encodedErr)
			assertMessageIs(t, grpcStatus, errorMessage)
		})
	}
}

func TestEncode_encodes_empty_message_for_OK(t *testing.T) {
	encodedErr := grpcerror.Encode(nil)

	grpcStatus := asGRPCError(t, encodedErr)

	const expected = ""
	if got := grpcStatus.Message(); got != expected {
		t.Errorf(
			`expected "%v", got "%v"`,
			expected,
			got,
		)
	}
}

func TestDecode_preserves_kind(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		t.Run(kind.String(), func(t *testing.T) {
			originalErr := apperror.New(kind, errorMessage)
			encodedErr := grpcerror.Encode(originalErr)
			decodedErr := grpcerror.Decode(encodedErr)

			got := apperror.KindOf(decodedErr)
			if got != kind {
				t.Errorf("expected %v, got %v", kind, got)
			}
		})
	}
}

func TestDecode_preserves_message_for_errors(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		if kind == apperror.OK {
			// nil errors do not have messages.
			continue
		}

		t.Run(kind.String(), func(t *testing.T) {
			originalErr := apperror.New(kind, errorMessage)
			encodedErr := grpcerror.Encode(originalErr)
			decodedErr := grpcerror.Decode(encodedErr)

			got := apperror.KindOf(decodedErr)
			if got != kind {
				t.Errorf("expected %v, got %v", kind, got)
			}
		})
	}
}

func TestDecode_preserves_nil_for_OK(t *testing.T) {
	originalErr := apperror.New(apperror.OK, errorMessage)
	encodedErr := grpcerror.Encode(originalErr)
	decodedErr := grpcerror.Decode(encodedErr)

	if decodedErr != nil {
		t.Error("expected nil, got", decodedErr)
	}
}

func TestDecode_return_unknown_error_on_failure(t *testing.T) {
	err := errors.New("not a gRPC error -- decoding will fail")
	decodedErr := grpcerror.Decode(err)

	const expected = apperror.UnknownError
	if got := apperror.KindOf(decodedErr); got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestEncodeKind_encoding_is_reversible(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		t.Run(kind.String(), func(t *testing.T) {
			code := grpcerror.EncodeKind(kind)
			got := grpcerror.DecodeKind(code)
			if got != kind {
				t.Errorf(
					"expected %v, but got %v",
					kind,
					got,
				)
			}
		})
	}
}

func TestDecodeKind_returns_unknown_for_codes_without_kind(t *testing.T) {
	for _, code := range []codes.Code{
		codes.Canceled,
		codes.Unknown,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unimplemented,
		codes.Internal,
		codes.DataLoss,
	} {
		t.Run(code.String(), func(t *testing.T) {
			got := grpcerror.DecodeKind(code)
			if expected := apperror.UnknownError; expected != got {
				t.Errorf("expected %v, got %v", expected, got)
			}
		})
	}
}

func asGRPCError(t *testing.T, err error) *status.Status {
	t.Helper()

	grpcStatus, ok := status.FromError(err)
	if !ok {
		t.Fatal("cannot parse gRPC error:", err)
	}

	return grpcStatus
}

func assertCodeIs(t *testing.T, s *status.Status, expected codes.Code) {
	t.Helper()

	if got := s.Code(); got != expected {
		t.Errorf("expected code %v, got %v", expected, got)
	}
}

func assertMessageIs(t *testing.T, s *status.Status, expected string) {
	t.Helper()

	if got := s.Message(); got != expected {
		t.Errorf(`expected message "%v", got "%v"`, expected, got)
	}
}

const errorMessage = "test error"
