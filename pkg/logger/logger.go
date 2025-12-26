package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// InitLogger initializes the global logger
func InitLogger(env string) {
	var handler slog.Handler
	if env == "development" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	Log = slog.New(handler)
	slog.SetDefault(Log)
}