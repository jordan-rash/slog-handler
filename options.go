package shandler

import (
	"io"
	"log/slog"
)

func WithJSON() HandlerOption {
	return func(h *Handler) {
		h.json = true
	}
}

func WithPid() HandlerOption {
	return func(h *Handler) {
		h.pid = true
	}
}

func WithStdOut(out ...io.Writer) HandlerOption {
	return func(h *Handler) {
		h.out = out
	}
}

func WithStdErr(err ...io.Writer) HandlerOption {
	return func(h *Handler) {
		h.err = err
	}
}

func WithTimeFormat(format string) HandlerOption {
	return func(h *Handler) {
		h.timeFormat = format
	}
}

func WithLineInfo(short bool) HandlerOption {
	return func(h *Handler) {
		h.lineInfo = true
		h.lineInfoShort = short
	}
}

// WithTextOutputFormat sets the format for the group output.
// The order of the fields are:
// 1. Record Level (Debug, Info, Warn, Error)
// 2. Record Time
// 3. Record Message
//
// The default format is "[%s] %s - %s\n".
// If you want to rearrange the fields, you can use the indexes:
// "%[3]s %[1]s %[2]\n"
// This will output "{Record Message} {Record Level} {Record Time}\n"
//
// User must provide newline in format
func WithTextOutputFormat(format string) HandlerOption {
	return func(h *Handler) {
		h.textOutputFormat = format
	}
}

func WithGroupTextOutputFormat(format string) HandlerOption {
	return func(h *Handler) {
		h.groupTextOutputFormat = format
	}
}

// WithGroupRightJustify will right justify the group name if set
// If there is an error calculating the width, it will default to 80
//
// Will override WithGroupTextOutputFormat setting if set
func WithGroupRightJustify() HandlerOption {
	return func(h *Handler) {
		h.groupRightJustify = true
	}
}

func WithLogLevel(level slog.Level) HandlerOption {
	return func(h *Handler) {
		h.level = level
	}
}

func WithColor() HandlerOption {
	return func(h *Handler) {
		h.color = true
	}
}

func WithTraceColor(color string) HandlerOption {
	return func(h *Handler) {
		h.traceColor = color
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

func WithFatalColor(color string) HandlerOption {
	return func(h *Handler) {
		h.fatalColor = color
	}
}

func WithShortLevels() HandlerOption {
	return func(h *Handler) {
		h.shortLevels = true
	}
}

// You can use this to filter out logs from specific groups
func WithGroupFilter(filter []string) HandlerOption {
	return func(h *Handler) {
		h.groupFilter = filter
	}
}
