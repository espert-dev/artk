package eventlog

import (
	"artk.dev/event"
	"context"
	"fmt"
	"log/slog"
)

func Logger[Event any](
	logger *slog.Logger,
) func(event.Observer[Event]) event.Observer[Event] {
	// Pre-generate the event type at initialization time, avoiding its
	// re-computation on every event.
	var example Event
	eventType := fmt.Sprintf("%T", example)
	logger = logger.With(slog.String(eventTypeKey, eventType))

	return func(next event.Observer[Event]) event.Observer[Event] {
		return func(ctx context.Context, e Event) error {
			logger = logger.With(slog.Any(eventKey, e))
			logger.LogAttrs(
				ctx,
				slog.LevelDebug,
				"event handler: notifying",
			)

			err := next(ctx, e)
			if err == nil {
				logger.LogAttrs(
					ctx,
					slog.LevelDebug,
					"event handler: success",
				)
			} else {
				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"event handler: failure",
					slog.String(errorKey, err.Error()),
				)
			}

			// Middleware must propagate errors.
			return err
		}
	}
}

const eventTypeKey = "eventType"
const eventKey = "event"
const errorKey = "error"
