package utils

import (
	"log/slog"
)

func WrapErr(err error) slog.Attr {
	return slog.String("error", err.Error())
}
