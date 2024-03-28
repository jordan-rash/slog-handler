package handler_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"
	"time"

	handler "github.com/jordan-rash/slog-handler"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			tt.opts = append(tt.opts, handler.WithStdOut(&stdout))
			logger := slog.New(handler.NewHandler(tt.opts...))

			logger.Error(tt.log)
			logger.Warn(tt.log)
			logger.Info(tt.log)
			logger.Debug(tt.log)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			tt.opts = append(tt.opts, handler.WithStdOut(&stdout), handler.WithJSON())
			logger := slog.New(handler.NewHandler(tt.opts...))

			logger.Error(tt.log)
			logger.Warn(tt.log)
			logger.Info(tt.log)
			logger.Debug(tt.log)
			assert.Equal(t, tt.expected, stdout.String())
		})
	}

}
