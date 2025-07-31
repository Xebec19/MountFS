package logger

import "log/slog"

var Logger *slog.Logger

func init() {
	Logger = slog.Default()
}

func Info(args ...any) {
	Logger.Info("log", args...)
}

func Error(args ...any) {
	Logger.Info("error", args...)
}
