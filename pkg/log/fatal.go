package log

import (
	"context"
	"log/slog"
	"os"
)

func Fatal(msg string, args ...any) {
	slog.Default().Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}
