package log

import "golang.org/x/exp/slog"

type message struct {
	Message string
	Level   slog.Level
	Attrs   []slog.Attr
}
