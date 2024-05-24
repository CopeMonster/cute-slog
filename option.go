package cute_slog

import "log/slog"

type Options struct {
	Level      slog.Level
	JsonFormat bool
	TimeFormat string
	SlotOpts   *slog.HandlerOptions
}
