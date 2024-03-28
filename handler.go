package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Handler struct {
	json             bool
	out              io.Writer
	err              io.Writer
	timeFormat       string
	textOutputFormat string
	level            slog.Level

	group string
	attrs []slog.Attr
}

type HandlerOption func(*Handler)

func NewHandler(opts ...HandlerOption) *Handler {
	nh := &Handler{
		attrs:            []slog.Attr{},
		out:              os.Stdout,
		err:              os.Stderr,
		timeFormat:       time.TimeOnly,
		textOutputFormat: "[%s] %s - %s\n",
		level:            slog.LevelInfo,
	}

	for _, opt := range opts {
		opt(nh)
	}

	return nh
}

func (n *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= n.level
}

func (n *Handler) Handle(ctx context.Context, record slog.Record) error {
	attrs := n.attrs
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	if !n.json {
		if record.NumAttrs() == 0 {
			fmt.Fprintf(n.out, n.textOutputFormat, record.Level, record.Time.Format(n.timeFormat), record.Message)
		} else {
			attsString := strings.Builder{}
			for _, a := range attrs {
				attsString.WriteString(a.String() + " ")
			}
			output := strings.TrimSpace(n.textOutputFormat) + " " + attsString.String() + "\n"
			fmt.Fprintf(n.out, output, record.Level, record.Time.Format(n.timeFormat), record.Message)
		}
	} else {
		a_map := make(map[string]any)
		for _, a := range attrs {
			a_map[a.Key] = a.Value.Any()
		}

		l := jsonLog{
			Level:   record.Level.String(),
			Time:    record.Time.Format(n.timeFormat),
			Message: record.Message,
			Group:   n.group,
			Attrs:   a_map,
		}

		l_raw, _ := json.Marshal(l)
		fmt.Fprintln(n.out, string(l_raw))
	}
	return nil
}

func (n *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := *n
	newHandler.attrs = append(newHandler.attrs, attrs...)
	return &newHandler
}

func (n *Handler) WithGroup(name string) slog.Handler {
	newHandler := *n
	newHandler.group = name
	return &newHandler
}

type jsonLog struct {
	Level   string         `json:"level"`
	Time    string         `json:"time"`
	Message string         `json:"message"`
	Group   string         `json:"group,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
}
