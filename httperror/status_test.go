package httperror_test

import (
	"artk.dev/apperror"
	"artk.dev/httperror"
	"net/http"
	"strconv"
	"testing"
)

func TestEncodeKind_decoding_reverses_encoding(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		t.Run(kind.String(), func(t *testing.T) {
			status := httperror.EncodeKind(kind)
			got := httperror.DecodeKind(status)
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

func TestDecodeKind_gateway_errors_are_timeouts(t *testing.T) {
	for _, code := range []int{
		http.StatusBadGateway,
		http.StatusGatewayTimeout,
	} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			got := httperror.DecodeKind(code)
			if expected := apperror.TimeoutError; expected != got {
				t.Errorf("expected %v, got %v", expected, got)
			}
		})
	}
}

func TestDecodeKind_misc_client_errors_are_validation_errors(t *testing.T) {
	// HTTP 405 is not one of the kinds considered by apperror because
	// it lacks meaning outside of an HTTP interface.
	got := httperror.DecodeKind(http.StatusMethodNotAllowed)
	if expected := apperror.ValidationError; expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestDecodeKind_misc_server_errors_are_validation_errors(t *testing.T) {
	// HTTP 505 is not one of the kinds considered by apperror because
	// it lacks meaning outside of an HTTP interface.
	got := httperror.DecodeKind(http.StatusHTTPVersionNotSupported)
	if expected := apperror.UnknownError; expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
