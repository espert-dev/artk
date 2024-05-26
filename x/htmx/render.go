package htmx

import (
	"context"
	"io"
	"net/http"
	"strings"
)

// Template invoked by Renderer.Render.
type Template interface {
	Render(ctx context.Context, w io.Writer) error
}

// Renderer renders an HTMX template to an http.ResponseWriter.
type Renderer struct {
	// OnError is synchronously called when rendering fails.
	OnError func(error)
}

// Render sets essential HTMX headers, renders the template, and notifies of
// any potential errors.
func (r Renderer) Render(
	ctx context.Context,
	w http.ResponseWriter,
	template Template,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Vary", varyHeaders)

	if err := template.Render(ctx, w); err != nil && r.OnError != nil {
		r.OnError(err)
	}
}

var varyHeaders = strings.Join([]string{
	acceptHeader,
	hxHistoryRestoreRequestHeader,
	hxRequestHeader,
	hxTriggerHeader,
}, ", ")
