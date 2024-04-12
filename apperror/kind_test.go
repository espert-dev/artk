package apperror_test

import (
	"errors"
	"github.com/jespert/artk/apperror"
	"strings"
	"testing"
)

// Test errors that pass the type checks but do not return their kind.

type timeoutErr struct {
	error
}

func (e timeoutErr) Timeout() bool {
	return true
}

type notModifiedErr struct {
	error
}

func (e notModifiedErr) NotModified() bool {
	return true
}

type validationErr struct {
	error
}

func (e validationErr) Validation() bool {
	return true
}

type unauthorizedErr struct {
	error
}

func (e unauthorizedErr) Unauthorized() bool {
	return true
}

type forbiddenErr struct {
	error
}

func (e forbiddenErr) Forbidden() bool {
	return true
}

type notFoundErr struct {
	error
}

func (e notFoundErr) NotFound() bool {
	return true
}

type conflictErr struct {
	error
}

func (e conflictErr) Conflict() bool {
	return true
}

type preconditionFailedErr struct {
	error
}

func (e preconditionFailedErr) PreconditionFailed() bool {
	return true
}

type tooManyRequestsErr struct {
	error
}

func (e tooManyRequestsErr) TooManyRequests() bool {
	return true
}

func TestKindOf(t *testing.T) {
	for _, tt := range []struct {
		kind apperror.Kind
		err  error
	}{
		{
			kind: apperror.UnknownKind,
			err:  errors.New("unknown error"),
		},
		{
			kind: apperror.NotModifiedKind,
			err:  notModifiedErr{},
		},
		{
			kind: apperror.ValidationKind,
			err:  validationErr{},
		},
		{
			kind: apperror.UnauthorizedKind,
			err:  unauthorizedErr{},
		},
		{
			kind: apperror.ForbiddenKind,
			err:  forbiddenErr{},
		},
		{
			kind: apperror.NotFoundKind,
			err:  notFoundErr{},
		},
		{
			kind: apperror.ConflictKind,
			err:  conflictErr{},
		},
		{
			kind: apperror.PreconditionFailedKind,
			err:  preconditionFailedErr{},
		},
		{
			kind: apperror.TooManyRequestsKind,
			err:  tooManyRequestsErr{},
		},
		{
			kind: apperror.TimeoutKind,
			err:  timeoutErr{},
		},
	} {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := apperror.KindOf(tt.err); got != tt.kind {
				t.Errorf(
					"expected %v, got %v",
					tt.kind,
					got,
				)
			}
		})
	}
}

func TestKind_String_contains_unexpected_values(t *testing.T) {
	if !strings.Contains(apperror.Kind(-1).String(), "-1") {
		t.Error("unexpected value (-1) not included")
	}
}
