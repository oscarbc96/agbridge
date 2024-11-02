package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

func newSlogZeroLogHandler(level slog.Level) slog.Handler {
	slogzerolog.ErrorKeys = []string{errKey}
	slogzerolog.LogLevels = map[slog.Level]zerolog.Level{
		LevelDebug: zerolog.DebugLevel,
		LevelInfo:  zerolog.InfoLevel,
		LevelWarn:  zerolog.WarnLevel,
		LevelError: zerolog.ErrorLevel,
		LevelFatal: zerolog.FatalLevel,
	}

	zerologLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly})

	return slogzerolog.Option{Level: level, Logger: &zerologLogger}.NewZerologHandler()
}
