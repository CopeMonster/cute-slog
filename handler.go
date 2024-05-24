package cute_slog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
)

type CuteHandler struct {
	Writer  io.Writer
	Options CuteOptions
}

func (c CuteHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= c.Options.Level
}

func (c CuteHandler) Handle(ctx context.Context, record slog.Record) error {
	var buf bytes.Buffer

	if c.Options.ColorFormat {
		switch record.Level {
		case slog.LevelInfo:
			buf.WriteString(GreenColor)
		case slog.LevelDebug:
			buf.WriteString(MagentaColor)
		case slog.LevelError:
			buf.WriteString(RedColor)
		case slog.LevelWarn:
			buf.WriteString(YellowColor)
		}

		buf.WriteString(BgBlackColor)
		buf.WriteString(fmt.Sprintf("[%s] ", record.Level.String()))
		buf.WriteString("\033[0m")
		buf.WriteString(record.Message)
	}

	if c.Options.LogTime {

	}

	if c.Options.JSONFormat {

	}

	_, err := c.Writer.Write(buf.Bytes())

	return err
}

func (c CuteHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]GroupOrAttrs, len(attrs))

	for i, attr := range attrs {
		newAttrs[i] = GroupOrAttrs{attr: attr}
	}

	return &CuteHandler{
		Writer:  c.Writer,
		Options: c.Options,
	}

}

func (c CuteHandler) WithGroup(name string) slog.Handler {
	return &CuteHandler{
		Writer:  c.Writer,
		Options: c.Options,
	}
}
