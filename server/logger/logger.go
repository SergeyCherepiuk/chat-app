package logger

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

type LogMessage struct {
	Message string
	Level   slog.Level
	Attrs   []slog.Attr
}

var LogMessages chan LogMessage

func init() {
	LogMessages = make(chan LogMessage)
}

func newLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{AddSource: true},
		),
	)
}

func HandleLogs() {
	logger := newLogger()
	for {
		logMessage := <-LogMessages
		logger.LogAttrs(
			context.Background(),
			logMessage.Level,
			logMessage.Message,
			logMessage.Attrs...,
		)
	}
}
