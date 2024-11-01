package log

import (
	"log/slog"

	slogmulti "github.com/samber/slog-multi"
)

func Setup(level slog.Level) {
	sink := slogmulti.Fanout(
		newSlogZeroLogHandler(level),
	)

	logger := slog.New(sink)

	slog.SetDefault(logger)
}
