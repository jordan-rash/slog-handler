package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	handler "disorder.dev/shandler"
)

func main() {
	logger := slog.New(handler.NewHandler(handler.WithJSON(), handler.WithPid()))
	logger.Info("test info")
	logger.Debug("test debug")
	logger.Info("test info", slog.String("key", "value"))
	logger.WithGroup("mygroup").Info("test info", slog.String("key", "value"))

	logger = slog.New(handler.NewHandler(handler.WithLogLevel(slog.LevelDebug), handler.WithColor(), handler.WithPid()))
	logger.Info("test info")
	logger.Debug("test debug")
	logger.Warn("test info", slog.String("key", "value"))
	logger.WithGroup("mygroup").Error("test info", slog.String("key", "value"))

	logger.With(slog.String("foo", "bar")).Info("test info", slog.String("key", "value"))
	logger.WithGroup("mygroup").Info("derp")

	logger = slog.New(handler.NewHandler(
		handler.WithLogLevel(slog.LevelDebug),
		handler.WithTimeFormat(time.RFC822),
		handler.WithTextOutputFormat("%s | %s | %s\n"),
		handler.WithStdErr(os.Stdout),
	))
	logger.With(slog.String("app", "myapp")).Debug("test")

	f, _ := os.Create("log.txt")
	defer f.Close()
	logger = slog.New(handler.NewHandler(
		handler.WithLogLevel(slog.LevelDebug),
		handler.WithStdOut(f),
		handler.WithStdErr(f),
	))
	logger.Info("test info")
	logger.Debug("test debug")
	err := errors.New("bad error")
	logger.Error("error", slog.Any("err", err))

	logger = slog.New(handler.NewHandler(
		handler.WithTextOutputFormat("%[3]s %[2]s %[1]s\n"),
	))
	logger.Info("flipped outout")

	logger = slog.New(handler.NewHandler(
		handler.WithLogLevel(handler.LevelTrace),
		handler.WithShortLevels(),
		handler.WithColor(),
	))
	logger.Info("testing trace next")
	logger.Log(context.Background(), handler.LevelTrace, "i am trace")

	logger = slog.New(handler.NewHandler(
		handler.WithLogLevel(slog.LevelError),
		handler.WithShortLevels(),
		handler.WithColor(),
	))
	logger.Error("error")
	logger.Log(context.Background(), handler.LevelFatal, "fatal")

	logger = slog.New(handler.NewHandler(
		handler.WithLogLevel(slog.LevelError),
		handler.WithPid(),
	))
	logger.Error("error")

	logger = slog.New(handler.NewHandler(
		handler.WithGroupRightJustify(),
	))
	logger.WithGroup("derp").Info("test")

	logger = slog.New(handler.NewHandler(
		handler.WithGroupRightJustify(),
	))
	logger.Info("test1")

	logger = slog.New(handler.NewHandler(
		handler.WithGroupRightJustify(),
		handler.WithPid(),
		handler.WithShortLevels(),
	))
	logger.Info("test1")
	logger.WithGroup("mygroup").Info("test2")
	logger.WithGroup("derp").Info("this is a long long long long long long message 1234567890")
}
