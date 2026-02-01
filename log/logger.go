package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func GetLogger(logLevel string) *slog.Logger {
	level := slog.LevelInfo
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      level,
		TimeFormat: time.Kitchen,
	}))
	return logger
}

func InitLogger(logLevel string) {
	slog.SetDefault(GetLogger(logLevel))
}
