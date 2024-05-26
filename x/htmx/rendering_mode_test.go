package htmx_test

import (
	"artk.dev/x/htmx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderingModeFor(t *testing.T) {
	for _, tc := range []struct {
		name     string
		headers  http.Header
		expected htmx.RenderingMode
	}{
		{
			name: "Prevent partials at top-level on navigation",
			headers: http.Header{
				"Hx-History-Restore-Request": []string{"true"},
				"Hx-Request":                 []string{"true"},
			},
			expected: htmx.FullPage,
		},
		{
			name: "Render partial on HTMX requests",
			headers: http.Header{
				"Hx-Request": []string{"true"},
			},
			expected: htmx.PartialUpdate,
		},
		{
			name:     "Render full by default",
			headers:  http.Header{},
			expected: htmx.FullPage,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header = tc.headers
			actual := htmx.RenderingModeFor(r)
			if tc.expected != actual {
				t.Errorf(
					"expected %v, got %v",
					tc.expected,
					actual,
				)
			}
		})
	}
}
