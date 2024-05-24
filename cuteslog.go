package cute_slog

import (
	"io"
	"log/slog"
)

func NewLogger(writer io.Writer, option CuteOptions) *slog.Logger {
	return slog.New(CuteHandler{
		Writer:  writer,
		Options: option,
	})
}
