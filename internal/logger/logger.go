package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var w = os.Stderr
var logger = slog.New(
	tint.NewHandler(w, &tint.Options{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	}),
)

func Default() *slog.Logger {
	return logger
}
