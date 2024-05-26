package htmx_test

import (
	"artk.dev/x/htmx"
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderer_Render_adds_expected_headers(t *testing.T) {
	var template Template
	template.RenderFn = func(_ context.Context, w io.Writer) error {
		_, err := w.Write([]byte("<html></html>"))
		return err
	}

	w := httptest.NewRecorder()
	r := htmx.Renderer{}
	r.Render(context.TODO(), w, template)
	response := w.Result()

	for _, tc := range []struct {
		header   string
		expected string
	}{
		{
			header:   "Content-Type",
			expected: "text/html; charset=utf-8",
		},
		{
			header: "Vary",
			expected: strings.Join([]string{
				"Accept",
				"Hx-History-Restore-Request",
				"Hx-Request",
				"Hx-Trigger",
			}, ", "),
		},
		{
			header:   "X-Content-Type-Options",
			expected: "nosniff",
		},
	} {
		t.Run(tc.header, func(t *testing.T) {
			actual := response.Header.Get(tc.header)
			if tc.expected != actual {
				t.Errorf(
					"expected %q, got %q",
					tc.expected,
					actual,
				)
			}
		})
	}
}

func TestRenderer_Render_supports_error_notifications(t *testing.T) {
	const expectedMsg = "test error"

	var template Template
	template.RenderFn = func(_ context.Context, _ io.Writer) error {
		return errors.New(expectedMsg)
	}

	var notifiedError error
	r := htmx.Renderer{
		OnError: func(err error) {
			notifiedError = err
		},
	}
	w := httptest.NewRecorder()
	r.Render(context.TODO(), w, template)

	if notifiedError == nil {
		t.Error("missing expected error")
	}
	if actualMsg := notifiedError.Error(); actualMsg != expectedMsg {
		t.Errorf("expected %q, got %q", expectedMsg, actualMsg)
	}
}

func TestRenderer_Render_does_not_panic_if_OnError_is_nil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("unexpected panic:", r)
		}
	}()

	var template Template
	template.RenderFn = func(_ context.Context, _ io.Writer) error {
		return errors.New("test error")
	}

	r := htmx.Renderer{}
	w := httptest.NewRecorder()
	r.Render(context.TODO(), w, template)
}

var _ htmx.Template = Template{}

type Template struct {
	RenderFn func(context.Context, io.Writer) error
}

func (t Template) Render(ctx context.Context, w io.Writer) error {
	return t.RenderFn(ctx, w)
}
