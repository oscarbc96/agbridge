package log

import (
	"errors"
)

func ParseLogLevel(levelStr string) (Level, error) {
	switch levelStr {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	case "fatal":
		return LevelFatal, nil
	default:
		return LevelInfo, errors.New("invalid log level: must be one of debug, info, warn, error, fatal")
	}
}
