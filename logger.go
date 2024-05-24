package cute_slog

import (
	"io"
	"log/slog"
)

type CLogger struct {
	*slog.Logger
}

func NewCLogger(writer io.Writer, options Options) CLogger {
	return CLogger{
		Logger: slog.New(NewHandler(writer, options)),
	}
}
