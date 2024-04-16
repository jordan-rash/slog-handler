package handler

import (
	"io"
	"log/slog"
)

func WithJSON() HandlerOption {
	return func(h *Handler) {
		h.json = true
	}
}

func WithStdOut(out io.Writer) HandlerOption {
	return func(h *Handler) {
		h.out = out
	}
}

func WithStdErr(err io.Writer) HandlerOption {
	return func(h *Handler) {
		h.err = err
	}
}

func WithTimeFormat(format string) HandlerOption {
	return func(h *Handler) {
		h.timeFormat = format
	}
}

func WithTextOutputFormat(format string) HandlerOption {
	return func(h *Handler) {
		h.textOutputFormat = format
	}
}

func WithLogLevel(level slog.Level) HandlerOption {
	return func(h *Handler) {
		h.level = level
	}
}

func WithColor(color bool) HandlerOption {
	return func(h *Handler) {
		h.color = color
	}
}

func WithDebugColor(color string) HandlerOption {
	return func(h *Handler) {
		h.debugColor = color
	}
}

func WithInfoColor(color string) HandlerOption {
	return func(h *Handler) {
		h.infoColor = color
	}
}

func WithWarnColor(color string) HandlerOption {
	return func(h *Handler) {
		h.warnColor = color
	}
}

func WithErrorColor(color string) HandlerOption {
	return func(h *Handler) {
		h.errorColor = color
	}
}
