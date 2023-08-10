package log

import "golang.org/x/exp/slog"

type Logger struct {
	attrs []slog.Attr
}

func (logger *Logger) With(attrs ...slog.Attr) {
	logger.attrs = append(logger.attrs, attrs...)
}

func (logger Logger) Info(msg string, attrs ...slog.Attr) {
	logs <- message{
		Message: msg,
		Level:   slog.LevelInfo,
		Attrs:   append(logger.attrs, attrs...),
	}
}

func (logger Logger) Warn(msg string, attrs ...slog.Attr) {
	logs <- message{
		Message: msg,
		Level:   slog.LevelWarn,
		Attrs:   append(logger.attrs, attrs...),
	}
}

func (logger Logger) Debug(msg string, attrs ...slog.Attr) {
	logs <- message{
		Message: msg,
		Level:   slog.LevelDebug,
		Attrs:   append(logger.attrs, attrs...),
	}
}

func (logger Logger) Error(msg string, attrs ...slog.Attr) {
	logs <- message{
		Message: msg,
		Level:   slog.LevelError,
		Attrs:   append(logger.attrs, attrs...),
	}
}