package log

import (
	"log/slog"
)

func Setup(level Level) {
	sink := newSlogZeroLogHandler(level)

	logger := slog.New(sink)

	slog.SetDefault(logger)
}
