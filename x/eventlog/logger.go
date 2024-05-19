package eventlog

import (
	"artk.dev/eventmux"
	"context"
	"fmt"
	"log/slog"
)

func Logger[Event any](
	logger *slog.Logger,
) func(eventmux.Observer[Event]) eventmux.Observer[Event] {
	// Pre-generate the event type at initialization time, avoiding its
	// re-computation on every event.
	var example Event
	eventType := fmt.Sprintf("%T", example)
	logger = logger.With(slog.String(eventTypeKey, eventType))

	return func(next eventmux.Observer[Event]) eventmux.Observer[Event] {
		return func(ctx context.Context, e Event) error {
			logger = logger.With(slog.Any(eventKey, e))
			logger.LogAttrs(
				ctx,
				slog.LevelDebug,
				"eventmux: notifying observer",
			)

			err := next(ctx, e)
			if err == nil {
				logger.LogAttrs(
					ctx,
					slog.LevelDebug,
					"eventmux: success",
				)
			} else {
				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"eventmux: failure",
					slog.String(errorKey, err.Error()),
				)
			}

			// Middleware must propagate errors.
			return err
		}
	}
}

const eventTypeKey = "event"
const eventKey = "eventType"
const errorKey = "error"
