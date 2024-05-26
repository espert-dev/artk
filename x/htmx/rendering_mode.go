//go:generate go run golang.org/x/tools/cmd/stringer@latest -output rendering_mode_string.go -type RenderingMode .
package htmx

import "net/http"

// RenderingMode indicates how to render the HTMX response.
type RenderingMode int

const (
	FullPage RenderingMode = iota
	PartialUpdate
)

// RenderingModeFor returns the RenderingMode for an http.Request.
func RenderingModeFor(r *http.Request) RenderingMode {
	header := r.Header
	if header.Get(hxHistoryRestoreRequestHeader) == "true" {
		return FullPage
	}
	if header.Get("Hx-Request") == "true" {
		return PartialUpdate
	}

	return FullPage
}
