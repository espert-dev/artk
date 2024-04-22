package testlog_test

import (
	"artk.dev/x/testlog"
	"context"
	"log/slog"
	"testing"
)

func TestNew(t *testing.T) {
	logger := testlog.New(t)
	logger.Log(
		context.TODO(),
		slog.LevelInfo,
		"Testing New function",
		slog.String("err", "boom"),
		slog.Int("answer", 42),
	)
}
