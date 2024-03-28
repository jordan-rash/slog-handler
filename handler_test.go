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

func TestTextLog(t *testing.T) {
	now := time.Now().Format(time.RFC822)

	var stdout bytes.Buffer
	logger := slog.New(handler.NewHandler(handler.WithStdOut(&stdout)))
	logger.Info("test")
	assert.Equal(t, fmt.Sprintf("[INFO] %s - test", now), stdout.String())
}
