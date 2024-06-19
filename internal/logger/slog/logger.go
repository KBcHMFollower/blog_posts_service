package sloglogger

import (
	"log/slog"
	"os"
)

const (
	loc  = "local"
	dev  = "development"
	prod = "production"
)

func ConfigureLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case loc, dev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
