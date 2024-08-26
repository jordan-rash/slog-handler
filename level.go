package shandler

import "log/slog"

const (
	LevelTrace slog.Level = slog.LevelDebug - 2
	LevelFatal slog.Level = slog.LevelError + 2
)
