package log

import (
	"log/slog"
)

const errKey = "err"

func Err(err error) slog.Attr {
	return slog.Any(errKey, err)
}
