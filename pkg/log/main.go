package log

import (
	"log/slog"
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelFatal = slog.Level(12)
)

var (
	With = slog.With

	Info  = slog.Info
	Debug = slog.Debug

	Any      = slog.Any
	Duration = slog.Duration
	Int      = slog.Int
	String   = slog.String
)

type Level = slog.Level

type Logger = slog.Logger
