package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
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
		timeFormat:       time.RFC822,
		textOutputFormat: "[%s] %s - %s",
		level:            slog.LevelInfo,
	}

	for _, opt := range opts {
		opt(nh)
	}

	return nh
}

func (n *Handler) Enabled(_ context.Context, level slog.Level) bool {
	fmt.Printf("level enabled: %t", level >= n.level)
	return level >= n.level
}

func (n *Handler) Handle(ctx context.Context, record slog.Record) error {
	fmt.Fprintf(n.out, n.textOutputFormat, record.Level, record.Time.Format(n.timeFormat), record.Message)
	return errors.New("invalid log level")
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
