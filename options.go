package cute_slog

import "log/slog"

type CuteOptions struct {
	Level       slog.Level
	ColorFormat bool
	JSONFormat  bool
	LogTime     bool
	TimeFormat  string
}

type GroupOrAttrs struct {
	attr  slog.Attr
	group string
}
