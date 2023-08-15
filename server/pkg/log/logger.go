package log

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

var logs chan message

func init() {
	logs = make(chan message, 10)
	for i := 0; i < 10; i++ {
		go handleLogs()
	}
}

func newLogger() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{},
		),
	)
}

func handleLogs() {
	logger := newLogger()
	for {
		logMessage := <-logs
		logger.LogAttrs(
			context.Background(),
			logMessage.Level,
			logMessage.Message,
			logMessage.Attrs...,
		)
	}
}
