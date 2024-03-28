package main

import (
	"log/slog"
	"os"
	"time"

	handler "github.com/jordan-rash/slog-handler"
)

func main() {
	logger := slog.New(handler.NewHandler(handler.WithJSON()))
	logger.Info("test info")
	logger.Debug("test debug")
	logger.Info("test info", slog.String("key", "value"))
	logger.WithGroup("mygroup").Info("test info", slog.String("key", "value"))

	logger = slog.New(handler.NewHandler(handler.WithLogLevel(slog.LevelDebug)))
	logger.Info("test info")
	logger.Debug("test debug")

	logger.With(slog.String("foo", "bar")).Info("test info", slog.String("key", "value"))
	logger.WithGroup("mygroup").Info("derp")

	logger = slog.New(handler.NewHandler(
		handler.WithLogLevel(slog.LevelDebug),
		handler.WithTimeFormat(time.RFC822),
		handler.WithTextOutputFormat("%s | %s | %s\n"),
		handler.WithStdErr(os.Stdout),
	))
	logger.With(slog.String("app", "myapp")).Debug("test")
}
