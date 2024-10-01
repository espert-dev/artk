package apperror_test

import (
	"artk.dev/apperror"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func Example_retrying_temporary_failures() {
	// An example operation that can fail temporarily.
	var i int
	operation := func() (int, error) {
		i++
		if i < 3 {
			return i, apperror.Timeoutf("timed out (%v)", i)
		}

		return i, nil
	}

	// Retry after a temporary failure up to a limited number of times.
	var value int
	var err error
	for range 3 {
		value, err = operation()

		// We went with IsFinal instead of IsTransient or IsTemporary
		// because most of the time you want to break out of a loop.
		// Since IsFinal returns true in that case, we avoid the NOT
		// operator and keep the logic positive. It is also shorter.
		if apperror.IsFinal(err) {
			break
		}

		fmt.Println("Operation failed:", err)
		time.Sleep(time.Nanosecond)
	}
	if err != nil {
		fmt.Println("Unexpected error:", err)
		return
	}

	fmt.Println("Success:", value)
	fmt.Println("Done.")

	// Output:
	// Operation failed: timed out (1)
	// Operation failed: timed out (2)
	// Success: 3
	// Done.
}

func Example_ignoring_an_already_deleted_item() {
	deleteItem := func(x int) error {
		if x == 2 {
			return apperror.NotFoundf("not found: %v", x)
		}

		return nil
	}

	// In this example, it is acceptable if the items are not found.
	for _, x := range []int{0, 2, 7, 42} {
		err := deleteItem(x)
		if apperror.IsNotFound(err) {
			fmt.Println("Already gone:", x)
			continue
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Deleted:", x)
	}

	fmt.Println("Done.")

	// Output:
	// Deleted: 0
	// Already gone: 2
	// Deleted: 7
	// Deleted: 42
	// Done.
}

func Example_checking_multiple_error_kinds_with_switch_only() {
	operation := func() error {
		return nil
	}

	switch err := operation(); apperror.KindOf(err) {
	case apperror.OK:
		// You can either handle the error here or leave it blank
		// and handle the normal path without indentation. In either
		// case, you likely want to have a case for OK so that it is
		// not handled by the default clause.
	case apperror.NotFoundError:
		fmt.Println("Not found:", err)
		return
	case apperror.ConflictError:
		fmt.Println("Conflict:", err)
		return
	case apperror.ForbiddenError:
		fmt.Println("Forbidden:", err)
		return
	default:
		fmt.Println("Other error:", err)
		return
	}

	// Handle the happy path.
	fmt.Println("OK.")

	// Output: OK.
}

func Example_checking_multiple_error_kinds_with_if_and_switch() {
	operation := func() error {
		return nil
	}

	// Nesting the switch inside a branch to separate happy and sad paths
	// could be of use sometimes.
	if err := operation(); err != nil {
		// Handle the sad path.
		switch apperror.KindOf(err) {
		case apperror.NotFoundError:
			fmt.Println("Not found:", err)
		case apperror.ConflictError:
			fmt.Println("Conflict:", err)
		case apperror.ForbiddenError:
			fmt.Println("Forbidden:", err)
		default:
			fmt.Println("Other error:", err)
		}
		return
	}

	// Handle the happy path.
	fmt.Println("OK.")

	// Output: OK.
}

type testCase struct {
	kind              apperror.Kind
	stringConstructor func(string) error
	formatConstructor func(string, ...any) error
	wrapper           func(error) error
	matcher           func(error) bool
}

func (tc testCase) Name() string {
	return tc.kind.String()
}

var testCases = []testCase{
	{
		kind:              apperror.ValidationError,
		stringConstructor: apperror.Validation,
		formatConstructor: apperror.Validationf,
		wrapper:           apperror.AsValidation,
		matcher:           apperror.IsValidation,
	},
	{
		kind:              apperror.UnauthorizedError,
		stringConstructor: apperror.Unauthorized,
		formatConstructor: apperror.Unauthorizedf,
		wrapper:           apperror.AsUnauthorized,
		matcher:           apperror.IsUnauthorized,
	},
	{
		kind:              apperror.ForbiddenError,
		stringConstructor: apperror.Forbidden,
		formatConstructor: apperror.Forbiddenf,
		wrapper:           apperror.AsForbidden,
		matcher:           apperror.IsForbidden,
	},
	{
		kind:              apperror.NotFoundError,
		stringConstructor: apperror.NotFound,
		formatConstructor: apperror.NotFoundf,
		wrapper:           apperror.AsNotFound,
		matcher:           apperror.IsNotFound,
	},
	{
		kind:              apperror.ConflictError,
		stringConstructor: apperror.Conflict,
		formatConstructor: apperror.Conflictf,
		wrapper:           apperror.AsConflict,
		matcher:           apperror.IsConflict,
	},
	{
		kind:              apperror.PreconditionFailedError,
		stringConstructor: apperror.PreconditionFailed,
		formatConstructor: apperror.PreconditionFailedf,
		wrapper:           apperror.AsPreconditionFailed,
		matcher:           apperror.IsPreconditionFailed,
	},
	{
		kind:              apperror.TooManyRequestsError,
		stringConstructor: apperror.TooManyRequests,
		formatConstructor: apperror.TooManyRequestsf,
		wrapper:           apperror.AsTooManyRequests,
		matcher:           apperror.IsTooManyRequests,
	},
	{
		kind:              apperror.UnknownError,
		stringConstructor: apperror.Unknown,
		formatConstructor: apperror.Unknownf,
		wrapper:           apperror.AsUnknown,
		matcher:           apperror.IsUnknown,
	},
	{
		kind:              apperror.TimeoutError,
		stringConstructor: apperror.Timeout,
		formatConstructor: apperror.Timeoutf,
		wrapper:           apperror.AsTimeout,
		matcher:           apperror.IsTimeout,
	},
}

func TestAs_accepts_empty_messages(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := apperror.As(tc.kind, errors.New(""))
			assertEmptyMessage(t, err)
		})
	}
}

func TestAs_invalid_kinds_are_mapped_to_unknown(t *testing.T) {
	const expected = apperror.UnknownError
	for _, kind := range []apperror.Kind{-1000000, 1000000} {
		err := apperror.As(kind, errors.New(""))
		if got := apperror.KindOf(err); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}

// This property can remove some unnecessary branches in client code.
//
// Consider:
//
//	return value, apperror.As(kind, err)
//
// vs.
//
//	if err == nil {
//	        return value, nil
//	} else {
//	        return nil, apperror.As(kind, err)
//	}
func TestAs_returns_nil_for_nil(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := apperror.As(tc.kind, nil)
			if err != nil {
				t.Error("expected nil, got:", err)
			}
		})
	}
}

func TestAs_returns_nil_for_OK(t *testing.T) {
	err := apperror.As(apperror.OK, errors.New(message))
	if err != nil {
		t.Error("expected nil, got:", err)
	}
}

// This property can eliminate some unnecessary branches in client code.
//
// Consider:
//
//	return value, apperror.AsValidation(err)
//
// vs.
//
//	if err == nil {
//	        return value, nil
//	} else {
//	        return nil, apperror.AsValidation(err)
//	}
func Test_wrappers_return_nil_for_nil(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.wrapper(nil)
			if err != nil {
				t.Error("expected nil, got:", err)
			}
		})
	}
}

func Test_matchers_return_false_for_nil(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			if tc.matcher(nil) {
				t.Error("matcher returned true for nil")
			}
		})
	}
}

func TestNew_accepts_empty_messages(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := apperror.New(tc.kind, "")
			assertEmptyMessage(t, err)
		})
	}
}

func TestNew_invalid_kinds_are_mapped_to_unknown(t *testing.T) {
	const expected = apperror.UnknownError
	for _, kind := range []apperror.Kind{-1000000, 1000000} {
		err := apperror.New(kind, message)
		if got := apperror.KindOf(err); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}

func TestNew_honors_kind(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := apperror.New(tc.kind, message)
			assertErrorKind(t, err, tc.kind, tc.matcher)
		})
	}
}

func TestNew_honors_message(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := apperror.Newf(tc.kind, "%v error", "test")
			assertTestMessage(t, err)
		})
	}
}

func TestNew_returns_nil_for_OK(t *testing.T) {
	err := apperror.Newf(apperror.OK, "%v error", "test")
	if err != nil {
		t.Error("expected nil, got:", err)
	}
}

func Test_string_constructors_accept_empty_messages(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.stringConstructor("")
			assertEmptyMessage(t, err)
		})
	}
}

func Test_format_constructors_accept_empty_messages(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.formatConstructor("")
			assertEmptyMessage(t, err)
		})
	}
}

func Test_string_constructors_honor_kind(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.stringConstructor(message)
			assertErrorKind(t, err, tc.kind, tc.matcher)
		})
	}
}

func Test_format_constructors_honor_kind(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.formatConstructor(message)
			assertErrorKind(t, err, tc.kind, tc.matcher)
		})
	}
}

func Test_string_constructors_honor_message(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.stringConstructor("test error")
			assertErrorKind(t, err, tc.kind, tc.matcher)
			assertTestMessage(t, err)
		})
	}
}

func Test_format_constructors_honor_message(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.formatConstructor("%v error", "test")
			assertErrorKind(t, err, tc.kind, tc.matcher)
			assertTestMessage(t, err)
		})
	}
}

func Test_wrappers_accept_empty_messages(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.wrapper(errors.New(""))
			assertEmptyMessage(t, err)
		})
	}
}

func Test_wrappers_honor_kind(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.wrapper(errors.New(message))
			assertErrorKind(t, err, tc.kind, tc.matcher)
		})
	}
}

func Test_wrappers_preserve_message(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name(), func(t *testing.T) {
			err := tc.wrapper(errors.New(message))
			assertTestMessage(t, err)
		})
	}
}

func Test_deadline_exceeded_is_a_timeout(t *testing.T) {
	assertErrorKind(
		t,
		context.DeadlineExceeded,
		apperror.TimeoutError,
		apperror.IsTimeout,
	)
}

func Test_non_semantic_errors_are_unknown(t *testing.T) {
	assertErrorKind(
		t,
		errors.New(message),
		apperror.UnknownError,
		apperror.IsUnknown,
	)
}

func TestIsUser(t *testing.T) {
	userKinds := map[apperror.Kind]struct{}{
		apperror.ValidationError:         {},
		apperror.UnauthorizedError:       {},
		apperror.ForbiddenError:          {},
		apperror.NotFoundError:           {},
		apperror.ConflictError:           {},
		apperror.PreconditionFailedError: {},
		apperror.TooManyRequestsError:    {},
	}

	for _, kind := range apperror.KindValues() {
		t.Run(kind.String(), func(t *testing.T) {
			_, expected := userKinds[kind]

			err := apperror.New(kind, message)
			got := apperror.IsUser(err)
			if got != expected {
				t.Errorf("expected %v, got %v", expected, got)
			}
		})
	}
}

func assertErrorKind(
	t *testing.T,
	err error,
	kind apperror.Kind,
	isExpectedKind func(error) bool,
) {
	t.Helper()

	if got := apperror.KindOf(err); got != kind {
		t.Errorf("unexpected kind %v, got %v", kind, got)
	}
	if !isExpectedKind(err) {
		t.Error("kind predicate failed")
	}
}

func assertTestMessage(t *testing.T, err error) {
	t.Helper()

	if got := err.Error(); got != message {
		t.Errorf(`Expected "%v", got "%v"`, message, got)
	}
}

func assertEmptyMessage(t *testing.T, err error) {
	t.Helper()

	if got := err.Error(); got != "" {
		t.Error("Expected empty error message, got:", got)
	}
}

const message = "test error"
