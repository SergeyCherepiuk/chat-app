package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

var Logger *slog.Logger

func init() {
	Logger = NewLogger()
}

func NewLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{AddSource: true},
		),
	)
}
