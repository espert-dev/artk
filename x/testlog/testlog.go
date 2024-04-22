// Package testlog provides colorized structured logging for tests.
package testlog

import (
	"github.com/lmittmann/tint"
	"github.com/neilotoole/slogt"
	"io"
	"log/slog"
	"testing"
)

// New creates a slog.Logger that will write to t.
func New(t *testing.T) *slog.Logger {
	return slogt.New(t, slogt.Factory(func(w io.Writer) slog.Handler {
		return tint.NewHandler(w, &tint.Options{
			Level: slog.LevelDebug,
		})
	}))
}
