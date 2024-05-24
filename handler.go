package cute_slog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log/slog"
	"strings"
	"sync"
)

type Handler struct {
	Options Options
	Writer  io.Writer
	attrs   []slog.Attr
	Mutex   *sync.Mutex
	slog.Handler
}

func NewHandler(writer io.Writer, options Options) *Handler {
	return &Handler{
		Options: options,
		Writer:  writer,
		attrs:   make([]slog.Attr, 0),
		Mutex:   new(sync.Mutex),
		Handler: slog.NewJSONHandler(writer, options.SlotOpts),
	}
}

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Options.Level
}

func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	var buf bytes.Buffer

	level := handleLevel(&record)

	b, err := handleAttributes(&record, h.Options.JsonFormat)
	if err != nil {
		return err
	}

	timeStr := handleTimestamp(&record, h.Options.TimeFormat)
	message := handleMessage(&record)

	writeToBuffer(&buf, timeStr, level, message, b)

	if _, err := h.Writer.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write buffer: %w", err)
	}

	return nil
}

func writeToBuffer(buf *bytes.Buffer, timeStr, level, message string, b []byte) {
	fmt.Fprintf(buf, "%s\t\t%s\t\t", timeStr, level)
	fmt.Fprintf(buf, "\033[0m")
	fmt.Fprintf(buf, message)
	if len(b) > 0 {
		fmt.Fprintf(buf, "%s\t%s\n", "\t", string(b))
	}
	fmt.Fprintf(buf, "\n")
}

func handleTimestamp(record *slog.Record, timeFormat string) string {
	return color.WhiteString(fmt.Sprintf("[%s]", record.Time.Format(timeFormat)))
}

func handleLevel(record *slog.Record) string {
	level := fmt.Sprintf("[%s]", record.Level.String())

	switch record.Level {
	case slog.LevelInfo:
		level = color.New(color.FgHiGreen, color.BgBlack).Sprintf("%s", level)
	case slog.LevelDebug:
		level = color.New(color.FgHiMagenta, color.BgBlack).Sprintf("%s", level)
	case slog.LevelError:
		level = color.New(color.FgHiRed, color.BgBlack).Sprintf("%s", level)
	case slog.LevelWarn:
		level = color.New(color.FgHiYellow, color.BgBlack).Sprintf("%s", level)
	}

	return level
}

func handleMessage(record *slog.Record) string {
	return fmt.Sprintf("%s", color.New(color.FgHiWhite).Sprintf("%s", record.Message))
}

func handleAttributes(record *slog.Record, formatJson bool) ([]byte, error) {
	var fields map[string]interface{}
	var attrs []string

	record.Attrs(func(attr slog.Attr) bool {
		if formatJson {
			if fields == nil {
				fields = make(map[string]interface{})
			}
			fields[attr.Key] = attr.Value.Any()
		} else {
			attrs = append(attrs, fmt.Sprintf("%s: %v", attr.Key, attr.Value.Any()))
		}
		return true
	})

	if formatJson && len(fields) > 0 {
		return json.MarshalIndent(fields, "", " ")
	} else if len(attrs) > 0 {
		return []byte(strings.Join(attrs, ", ")), nil
	}

	return nil, nil
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		Options: h.Options,
		Writer:  h.Writer,
		attrs:   attrs,
		Mutex:   h.Mutex,
		Handler: h.Handler,
	}
}

func (h Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		Options: h.Options,
		Writer:  h.Writer,
		attrs:   make([]slog.Attr, 0),
		Mutex:   h.Mutex,
		Handler: h.Handler.WithGroup(name),
	}
}
