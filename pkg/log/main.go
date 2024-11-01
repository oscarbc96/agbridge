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
	With  = slog.With
	Group = slog.Group

	Debug = slog.Debug
	Info  = slog.Info
	Warn  = slog.Warn
	Error = slog.Error

	Any    = slog.Any
	String = slog.String
	Time   = slog.Time
	Int    = slog.Int
	Bool   = slog.Bool
)

type Logger = slog.Logger
