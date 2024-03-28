package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
	"time"

	handler "github.com/jordan-rash/slog-handler"
	"github.com/stretchr/testify/assert"
)

func TestTextLog(t *testing.T) {
	now := time.Now().Format(time.RFC822)

	var stdout bytes.Buffer
	logger := slog.New(handler.NewHandler(handler.WithStdOut(&stdout)))
	logger.Info("test")
	assert.Equal(t, fmt.Sprintf("[INFO] %s - test", now), stdout.String())
}

func TestJsonLog(t *testing.T) {
	now := time.Now().Format(time.RFC822)
	expected := struct {
		Level   string `json:"level"`
		Time    string `json:"time"`
		Message string `json:"message"`
	}{
		Level:   "INFO",
		Time:    now,
		Message: "test",
	}
	expected_r, _ := json.Marshal(expected)

	var stdout bytes.Buffer
	logger := slog.New(handler.NewHandler(handler.WithStdOut(&stdout), handler.WithJSON()))
	logger.Info("test")
	assert.Equal(t, string(expected_r), stdout.String())
}
