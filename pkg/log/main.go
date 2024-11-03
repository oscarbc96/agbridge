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
	Group = slog.Group
	With  = slog.With

	Debug = slog.Debug
	Error = slog.Error
	Info  = slog.Info
	Warn  = slog.Warn

	Any      = slog.Any
	Bool     = slog.Bool
	Duration = slog.Duration
	Int      = slog.Int
	String   = slog.String
	Time     = slog.Time
)

type Level = slog.Level

type Logger = slog.Logger
