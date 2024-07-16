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

	"github.com/charmbracelet/lipgloss"
)

type Handler struct {
	json                  bool
	shortLevels           bool
	out                   []io.Writer
	err                   []io.Writer
	timeFormat            string
	textOutputFormat      string
	groupTextOutputFormat string
	level                 slog.Level

	color      bool
	traceColor string
	debugColor string
	infoColor  string
	warnColor  string
	errorColor string

	group string
	attrs []slog.Attr
}

type HandlerOption func(*Handler)

func NewHandler(opts ...HandlerOption) *Handler {
	nh := &Handler{
		attrs:                 []slog.Attr{},
		out:                   []io.Writer{os.Stdout},
		err:                   []io.Writer{os.Stderr},
		shortLevels:           false,
		timeFormat:            time.TimeOnly,
		textOutputFormat:      "[%s] %s - %s\n",
		groupTextOutputFormat: "%s | %s",
		level:                 slog.LevelInfo,
		color:                 false,
		traceColor:            "#C0C0C0", // Gray
		debugColor:            "#FFE6FF", // Light Pink
		infoColor:             "#6666FF", // Slate Blue
		warnColor:             "#FFBB33", // Burnt Orange
		errorColor:            "#E60000", // Crimson Red
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

	textFormat := func() string {
		if n.group != "" {
			return fmt.Sprintf(n.groupTextOutputFormat, n.group, n.textOutputFormat)
		}
		return n.textOutputFormat
	}

	outLoc := func() []io.Writer {
		if record.Level >= slog.LevelError {
			return n.err
		}
		return n.out
	}

	level := func() string {
		if n.shortLevels {
			switch record.Level {
			case LevelTrace:
				return "TRC"
			case slog.LevelDebug:
				return "DBG"
			case slog.LevelInfo:
				return "INF"
			case slog.LevelWarn:
				return "WRN"
			case slog.LevelError:
				return "ERR"
			}
		} else {
			switch record.Level {
			case LevelTrace:
				return "TRACE"
			}
		}
		return record.Level.String()
	}

	var recordLevel string
	if n.color {
		switch record.Level {
		case LevelTrace:
			recordLevel = lipgloss.NewStyle().Foreground(lipgloss.Color(n.traceColor)).Render(level())
		case slog.LevelDebug:
			recordLevel = lipgloss.NewStyle().Foreground(lipgloss.Color(n.debugColor)).Render(level())
		case slog.LevelInfo:
			recordLevel = lipgloss.NewStyle().Foreground(lipgloss.Color(n.infoColor)).Render(level())
		case slog.LevelWarn:
			recordLevel = lipgloss.NewStyle().Foreground(lipgloss.Color(n.warnColor)).Render(level())
		case slog.LevelError:
			recordLevel = lipgloss.NewStyle().Foreground(lipgloss.Color(n.errorColor)).Render(level())
		}
	} else {
		recordLevel = level()
	}

	if !n.json {
		if len(attrs) == 0 {
			printerf(outLoc(), textFormat(), recordLevel, record.Time.Format(n.timeFormat), record.Message)
		} else {
			attsString := strings.Builder{}
			for i, a := range attrs {
				attsString.WriteString(a.String())
				if i < len(attrs)-1 {
					attsString.WriteString(" ")
				}
			}
			output := strings.TrimSpace(textFormat()) + " " + attsString.String() + "\n"
			printerf(outLoc(), output, recordLevel, record.Time.Format(n.timeFormat), record.Message)
		}
	} else {
		a_map := make(map[string]any)
		for _, a := range attrs {
			a_map[a.Key] = a.Value.Any()
		}

		l := jsonLog{
			Level:   level(),
			Time:    record.Time.Format(n.timeFormat),
			Message: record.Message,
			Group:   n.group,
			Attrs:   a_map,
		}

		l_raw, _ := json.Marshal(l)
		printer(outLoc(), string(l_raw))
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
