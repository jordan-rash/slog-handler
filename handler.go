package shandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var width int = 0

func init() {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
		return
	}
	width = w - 1
}

type Handler struct {
	json        bool
	pid         bool
	shortLevels bool

	lineInfo      bool
	lineInfoShort bool

	out                   []io.Writer
	err                   []io.Writer
	timeFormat            string
	textOutputFormat      string
	groupTextOutputFormat string
	groupRightJustify     bool
	level                 slog.Level

	color      bool
	traceColor string
	debugColor string
	infoColor  string
	warnColor  string
	errorColor string
	fatalColor string

	group       string
	groupFilter []string
	attrs       []slog.Attr
}

type HandlerOption func(*Handler)

func NewHandler(opts ...HandlerOption) *Handler {
	nh := &Handler{
		pid:                   false,
		attrs:                 []slog.Attr{},
		out:                   []io.Writer{os.Stdout},
		err:                   []io.Writer{os.Stderr},
		shortLevels:           false,
		lineInfo:              false,
		lineInfoShort:         true,
		timeFormat:            time.TimeOnly,
		textOutputFormat:      "[%s] %s - %s\n",
		groupTextOutputFormat: "%s | %s",
		groupFilter:           []string{},
		level:                 slog.LevelInfo,
		color:                 false,
		traceColor:            "#C0C0C0", // Gray
		debugColor:            "#FFE6FF", // Light Pink
		infoColor:             "#6666FF", // Slate Blue
		warnColor:             "#FFBB33", // Burnt Orange
		errorColor:            "#E60000", // Crimson Red
		fatalColor:            "#990000", // Dark Red
	}

	for _, opt := range opts {
		opt(nh)
	}

	return nh
}

// NewHandlerFromConfig will allow you to pass in the settings to
// slog.New(NewHandlerFromConfig). You will need to inclused the io.Writers
// in the NewHandlerFromConfig call as they are not serializable.
// Use ToConfig to get the config of your original Handler
func NewHandlerFromConfig(config []byte, stdout, stderr []io.Writer) (*Handler, error) {
	nh := new(Handler)
	if err := json.Unmarshal(config, nh); err != nil {
		return nil, err
	}

	nh.out = stdout
	nh.err = stderr
	return nh, nil
}

func (n *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= n.level
}

func (n *Handler) Handle(ctx context.Context, record slog.Record) error {
	if slices.Contains(n.groupFilter, n.group) {
		return nil
	}

	attrs := n.attrs
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	// This was adapted from stdlib record.go:219
	if n.lineInfo {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, _ := fs.Next()

		var logLine string
		if n.lineInfoShort {
			fileBase := filepath.Base(f.File)
			logLine = fmt.Sprintf("%s:%d", fileBase, f.Line)
		} else {
			logLine = fmt.Sprintf("%s [%s:%d]", f.Function, f.File, f.Line)
		}

		if logLine != "" {
			attrs = append(attrs, slog.String("slog_info", logLine))
		} else {
			attrs = append(attrs, slog.String("slog_info", "unknown"))
		}
	}

	textFormat := func() string {
		if n.group != "" && !n.groupRightJustify {
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
			case LevelFatal:
				return "FTL"
			}
		} else {
			switch record.Level {
			case LevelTrace:
				return "TRACE"
			case LevelFatal:
				return "FATAL"
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
		case LevelFatal:
			recordLevel = lipgloss.NewStyle().Foreground(lipgloss.Color(n.fatalColor)).Render(level())
		}
	} else {
		recordLevel = level()
	}

	var pid string
	if n.pid {
		pid = strconv.Itoa(os.Getpid())
	}

	if !n.json {
		output := textFormat()
		attsString := strings.Builder{}

		if len(attrs) != 0 {
			output = strings.TrimSpace(output)
			for i, a := range attrs {
				attsString.WriteString(a.String())
				if i < len(attrs)-1 {
					attsString.WriteString(" ")
				}
			}
			attsString.WriteString("\n")
			output = output + " " + attsString.String()
		}

		if n.groupRightJustify {
			printerrj(outLoc(), n.group, pid, output, recordLevel, record.Time.Format(n.timeFormat), record.Message)
		} else {
			printerf(outLoc(), pid, output, recordLevel, record.Time.Format(n.timeFormat), record.Message)
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
			Pid:     pid,
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

func ToConfig(h slog.Handler) ([]byte, error) {
	n, ok := (h).(*Handler)
	if !ok {
		return nil, fmt.Errorf("not a shandler")
	}
	return json.Marshal(n)
}

type attrValue struct {
	Kind  slog.Kind `json:"kind"`
	Value string    `json:"value"`
}

func (n Handler) MarshalJSON() ([]byte, error) {
	a := map[string]attrValue{}
	for _, aT := range n.attrs {
		a[aT.Key] = attrValue{Kind: aT.Value.Kind(), Value: aT.Value.String()}
	}

	return json.Marshal(map[string]any{
		"json":                     n.json,
		"short_levels":             n.shortLevels,
		"line_info":                n.lineInfo,
		"time_format":              n.timeFormat,
		"text_output_format":       n.textOutputFormat,
		"group_text_output_format": n.groupTextOutputFormat,
		"level":                    n.level,
		"color":                    n.color,
		"trace_color":              n.traceColor,
		"debug_color":              n.debugColor,
		"info_color":               n.infoColor,
		"warn_color":               n.warnColor,
		"error_color":              n.errorColor,
		"fatal_color":              n.fatalColor,
		"group":                    n.group,
		"group_filter":             n.groupFilter,
		"attrs":                    a,
	})
}

func (n *Handler) UnmarshalJSON(data []byte) error {
	temp := struct {
		Json                  bool                 `json:"json"`
		ShortLevels           bool                 `json:"short_levels"`
		LineInfo              bool                 `json:"line_info"`
		TimeFormat            string               `json:"time_format"`
		TextOutputFormat      string               `json:"text_output_format"`
		GroupTextOutputFormat string               `json:"group_text_output_format"`
		Level                 slog.Level           `json:"level"`
		Color                 bool                 `json:"color"`
		TraceColor            string               `json:"trace_color"`
		DebugColor            string               `json:"debug_color"`
		InfoColor             string               `json:"info_color"`
		WarnColor             string               `json:"warn_color"`
		ErrorColor            string               `json:"error_color"`
		FatalColor            string               `json:"fatal_color"`
		Group                 string               `json:"group"`
		GroupFilter           []string             `json:"group_filter"`
		Attrs                 map[string]attrValue `json:"attrs"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	n.json = temp.Json
	n.shortLevels = temp.ShortLevels
	n.timeFormat = temp.TimeFormat
	n.textOutputFormat = temp.TextOutputFormat
	n.groupTextOutputFormat = temp.GroupTextOutputFormat
	n.level = temp.Level
	n.color = temp.Color
	n.traceColor = temp.TraceColor
	n.debugColor = temp.DebugColor
	n.infoColor = temp.InfoColor
	n.warnColor = temp.WarnColor
	n.errorColor = temp.ErrorColor
	n.fatalColor = temp.FatalColor
	n.group = temp.Group
	n.groupFilter = temp.GroupFilter

	for k, v := range temp.Attrs {
		switch v.Kind {
		case slog.KindAny:
			n.attrs = append(n.attrs, slog.Any(k, v.Value))
		case slog.KindBool:
			n.attrs = append(n.attrs, slog.Bool(k, v.Value == "true"))
		case slog.KindDuration:
			d, _ := time.ParseDuration(v.Value)
			n.attrs = append(n.attrs, slog.Duration(k, d))
		case slog.KindFloat64:
			num, err := strconv.ParseFloat(v.Value, 64)
			if err != nil {
				return err
			}
			n.attrs = append(n.attrs, slog.Float64(k, num))
		case slog.KindInt64:
			num, err := strconv.ParseInt(v.Value, 10, 64)
			if err != nil {
				return err
			}
			n.attrs = append(n.attrs, slog.Int64(k, num))
		case slog.KindString:
			n.attrs = append(n.attrs, slog.String(k, v.Value))
		case slog.KindTime:
			t, err := time.Parse(n.timeFormat, v.Value)
			if err != nil {
				return err
			}
			n.attrs = append(n.attrs, slog.Time(k, t))
		case slog.KindUint64:
			num, err := strconv.ParseUint(v.Value, 10, 64)
			if err != nil {
				return err
			}
			n.attrs = append(n.attrs, slog.Uint64(k, num))
		}
	}
	return nil
}

type jsonLog struct {
	Level   string         `json:"level"`
	Time    string         `json:"time"`
	Message string         `json:"message"`
	Group   string         `json:"group,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Pid     string         `json:"pid,omitempty"`
}
