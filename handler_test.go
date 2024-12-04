package shandler_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	handler "disorder.dev/shandler"
	"github.com/stretchr/testify/assert"
)

func TestNewHandlerText(t *testing.T) {
	now := time.Now().Format(time.TimeOnly)

	tests := []struct {
		name     string
		opts     []handler.HandlerOption
		log      string
		expected string
	}{
		{name: "defaults", opts: []handler.HandlerOption{}, log: "test", expected: fmt.Sprintf("[ERROR] %s - test\n[WARN] %s - test\n[INFO] %s - test\n", now, now, now)},
		{name: "error_level", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelError)}, log: "test", expected: fmt.Sprintf("[ERROR] %s - test\n", now)},
		{name: "warn_level", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelWarn)}, log: "test", expected: fmt.Sprintf("[ERROR] %s - test\n[WARN] %s - test\n", now, now)},
		{name: "info_level", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelInfo)}, log: "test", expected: fmt.Sprintf("[ERROR] %s - test\n[WARN] %s - test\n[INFO] %s - test\n", now, now, now)},
		{name: "debug_level", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelDebug)}, log: "test", expected: fmt.Sprintf("[ERROR] %s - test\n[WARN] %s - test\n[INFO] %s - test\n[DEBUG] %s - test\n", now, now, now, now)},
		{name: "debug_level_shortlvl", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelDebug), handler.WithShortLevels()}, log: "test", expected: fmt.Sprintf("[ERR] %s - test\n[WRN] %s - test\n[INF] %s - test\n[DBG] %s - test\n", now, now, now, now)},
		{name: "debug_level_group", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelDebug)}, log: "test", expected: fmt.Sprintf("group | [ERROR] %s - test\ngroup | [WARN] %s - test\ngroup | [INFO] %s - test\ngroup | [DEBUG] %s - test\n", now, now, now, now)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			tt.opts = append(tt.opts, handler.WithStdOut(&stdout), handler.WithStdErr(&stdout))
			logger := slog.New(handler.NewHandler(tt.opts...))

			if strings.HasSuffix(tt.name, "_group") {
				logger.WithGroup("group").Error(tt.log)
				logger.WithGroup("group").Warn(tt.log)
				logger.WithGroup("group").Info(tt.log)
				logger.WithGroup("group").Debug(tt.log)
			} else {
				logger.Error(tt.log)
				logger.Warn(tt.log)
				logger.Info(tt.log)
				logger.Debug(tt.log)
			}
			assert.Equal(t, tt.expected, stdout.String())
		})
	}
}

func TestJsonLog(t *testing.T) {
	now := time.Now().Format(time.TimeOnly)
	tests := []struct {
		name     string
		opts     []handler.HandlerOption
		log      string
		expected string
	}{
		{name: "defaults", opts: []handler.HandlerOption{}, log: "test", expected: fmt.Sprintf("{\"level\":\"ERROR\",\"time\":\"%s\",\"message\":\"test\"}\n{\"level\":\"WARN\",\"time\":\"%s\",\"message\":\"test\"}\n{\"level\":\"INFO\",\"time\":\"%s\",\"message\":\"test\"}\n", now, now, now)},
		{name: "error_level", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelError)}, log: "test", expected: fmt.Sprintf("{\"level\":\"ERROR\",\"time\":\"%s\",\"message\":\"test\"}\n", now)},
		{name: "warn_level", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelWarn)}, log: "test", expected: fmt.Sprintf("{\"level\":\"ERROR\",\"time\":\"%s\",\"message\":\"test\"}\n{\"level\":\"WARN\",\"time\":\"%s\",\"message\":\"test\"}\n", now, now)},
		{name: "info_level", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelInfo)}, log: "test", expected: fmt.Sprintf("{\"level\":\"ERROR\",\"time\":\"%s\",\"message\":\"test\"}\n{\"level\":\"WARN\",\"time\":\"%s\",\"message\":\"test\"}\n{\"level\":\"INFO\",\"time\":\"%s\",\"message\":\"test\"}\n", now, now, now)},
		{name: "debug_level", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelDebug)}, log: "test", expected: fmt.Sprintf("{\"level\":\"ERROR\",\"time\":\"%s\",\"message\":\"test\"}\n{\"level\":\"WARN\",\"time\":\"%s\",\"message\":\"test\"}\n{\"level\":\"INFO\",\"time\":\"%s\",\"message\":\"test\"}\n{\"level\":\"DEBUG\",\"time\":\"%s\",\"message\":\"test\"}\n", now, now, now, now)},
		{name: "debug_level_group", opts: []handler.HandlerOption{handler.WithLogLevel(slog.LevelDebug)}, log: "test", expected: fmt.Sprintf("{\"level\":\"ERROR\",\"time\":\"%s\",\"message\":\"test\",\"group\":\"group\"}\n{\"level\":\"WARN\",\"time\":\"%s\",\"message\":\"test\",\"group\":\"group\"}\n{\"level\":\"INFO\",\"time\":\"%s\",\"message\":\"test\",\"group\":\"group\"}\n{\"level\":\"DEBUG\",\"time\":\"%s\",\"message\":\"test\",\"group\":\"group\"}\n", now, now, now, now)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			tt.opts = append(tt.opts, handler.WithStdOut(&stdout), handler.WithStdErr(&stdout), handler.WithJSON())
			logger := slog.New(handler.NewHandler(tt.opts...))

			if strings.HasSuffix(tt.name, "_group") {
				logger.WithGroup("group").Error(tt.log)
				logger.WithGroup("group").Warn(tt.log)
				logger.WithGroup("group").Info(tt.log)
				logger.WithGroup("group").Debug(tt.log)
			} else {
				logger.Error(tt.log)
				logger.Warn(tt.log)
				logger.Info(tt.log)
				logger.Debug(tt.log)
			}
			assert.Equal(t, tt.expected, stdout.String())
		})
	}
}

func TestWithJsonOption(t *testing.T) {
	var stdout bytes.Buffer
	now := time.Now().Format(time.TimeOnly)

	logger := slog.New(handler.NewHandler(handler.WithJSON(), handler.WithStdOut(&stdout)))
	logger.Info("test")

	assert.Equal(t, fmt.Sprintf("{\"level\":\"INFO\",\"time\":\"%s\",\"message\":\"test\"}\n", now), stdout.String())
}

func TestWithTimeFormatOption(t *testing.T) {
	var stdout bytes.Buffer
	now := time.Now().Format(time.RFC822)

	logger := slog.New(handler.NewHandler(handler.WithTimeFormat(time.RFC822), handler.WithStdOut(&stdout)))
	logger.Info("test")

	assert.Equal(t, fmt.Sprintf("[INFO] %s - test\n", now), stdout.String())
}

func TestWithTextOutputFormatOption(t *testing.T) {
	var stdout bytes.Buffer
	now := time.Now().Format(time.TimeOnly)

	logger := slog.New(handler.NewHandler(handler.WithTextOutputFormat("%s | %s | %s\n"), handler.WithStdOut(&stdout)))
	logger.Info("test")

	assert.Equal(t, fmt.Sprintf("INFO | %s | test\n", now), stdout.String())
}

func TestLogAttr(t *testing.T) {
	var stdout bytes.Buffer
	now := time.Now().Format(time.TimeOnly)

	logger := slog.New(handler.NewHandler(handler.WithStdErr(&stdout), handler.WithStdOut(&stdout)))
	logger.With(slog.String("foo", "bar")).Info("test")

	assert.Equal(t, fmt.Sprintf("[INFO] %s - test foo=bar\n", now), stdout.String())
}

func TestLoggerAttr(t *testing.T) {
	var stdout bytes.Buffer
	now := time.Now().Format(time.TimeOnly)

	logger := slog.New(handler.NewHandler(handler.WithStdOut(&stdout))).With(slog.String("foo", "bar"))
	logger.Info("test")

	assert.Equal(t, fmt.Sprintf("[INFO] %s - test foo=bar\n", now), stdout.String())
}

func TestGroupWithNewFormat(t *testing.T) {
	var stdout bytes.Buffer
	now := time.Now().Format(time.TimeOnly)

	logger := slog.New(handler.NewHandler(handler.WithStdOut(&stdout), handler.WithTextOutputFormat("[%s] %s - %s"), handler.WithGroupTextOutputFormat("%[2]s <> %[1]s\n"))).WithGroup("mygroup")
	logger.Info("test")

	assert.Equal(t, fmt.Sprintf("[INFO] %s - test <> mygroup\n", now), stdout.String())
}

func TestTraceLevel(t *testing.T) {
	var stdout bytes.Buffer
	now := time.Now().Format(time.TimeOnly)

	logger := slog.New(handler.NewHandler(handler.WithStdOut(&stdout), handler.WithLogLevel(handler.LevelTrace)))
	logger.Log(context.TODO(), handler.LevelTrace, "test")

	assert.Equal(t, fmt.Sprintf("[TRACE] %s - test\n", now), stdout.String())

	stdout = bytes.Buffer{}
	logger = slog.New(handler.NewHandler(handler.WithStdOut(&stdout), handler.WithLogLevel(handler.LevelTrace), handler.WithShortLevels()))
	logger.Log(context.TODO(), handler.LevelTrace, "test")

	assert.Equal(t, fmt.Sprintf("[TRC] %s - test\n", now), stdout.String())
}

func TestFatalLevel(t *testing.T) {
	var stderr bytes.Buffer
	now := time.Now().Format(time.TimeOnly)

	logger := slog.New(handler.NewHandler(handler.WithStdErr(&stderr), handler.WithLogLevel(handler.LevelFatal)))
	logger.Log(context.TODO(), handler.LevelFatal, "test")

	assert.Equal(t, fmt.Sprintf("[FATAL] %s - test\n", now), stderr.String())

	stderr = bytes.Buffer{}
	logger = slog.New(handler.NewHandler(handler.WithStdErr(&stderr), handler.WithLogLevel(handler.LevelFatal), handler.WithShortLevels()))
	logger.Log(context.TODO(), handler.LevelFatal, "test")

	assert.Equal(t, fmt.Sprintf("[FTL] %s - test\n", now), stderr.String())
}

func TestGroupFilter(t *testing.T) {
	var stdout bytes.Buffer
	now := time.Now().Format(time.TimeOnly)
	logger := slog.New(handler.NewHandler(handler.WithStdOut(&stdout), handler.WithGroupFilter([]string{"group"})))
	logger.Info("line 1")
	logger = logger.WithGroup("group")
	logger.Info("line 2")
	logger = logger.WithGroup("foo")
	logger.Info("line 3")
	assert.Equal(t, fmt.Sprintf("[INFO] %s - line 1\nfoo | [INFO] %s - line 3\n", now, now), stdout.String())
}

func TestFromConfig(t *testing.T) {
	var stdout bytes.Buffer
	now := time.Now().Format(time.TimeOnly)
	logger := slog.New(handler.NewHandler(handler.WithStdOut(&stdout))).WithGroup("copyme").With(slog.String("foo", "bar"))
	logger_config, err := handler.ToConfig(logger.Handler())
	assert.Nil(t, err)

	config_logger, err := handler.NewHandlerFromConfig(logger_config, []io.Writer{&stdout}, nil)
	assert.Nil(t, err)
	slog.New(config_logger).Info("test")
	assert.Equal(t, fmt.Sprintf("copyme | [INFO] %s - test foo=bar\n", now), stdout.String())
}

func BenchmarkHandlers(b *testing.B) {
	var stdout bytes.Buffer
	bt := []struct {
		Name    string
		Handler slog.Handler
	}{
		{"handler text log", handler.NewHandler(handler.WithStdOut(&stdout))},
		{"stdlib text log", slog.NewTextHandler(&stdout, nil)},
		{"handler json log", handler.NewHandler(handler.WithStdOut(&stdout), handler.WithJSON())},
		{"stdlib json log", slog.NewJSONHandler(&stdout, nil)},
	}

	for _, t := range bt {
		b.Run(t.Name, func(b *testing.B) {
			logger := slog.New(t.Handler)
			for i := 0; i < b.N; i++ {
				logger.Info("test")
			}
		})
		stdout = bytes.Buffer{}
	}
}
