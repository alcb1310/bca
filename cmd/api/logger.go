package main

import (
	"log/slog"
	"os"
	"strings"
)

func loggerSetup(env string) {
	handlerOptions := &slog.HandlerOptions{}

	switch strings.ToLower(env) {
	case "debug":
		handlerOptions.Level = slog.LevelDebug
	case "info":
		handlerOptions.Level = slog.LevelInfo
	case "warn":
		handlerOptions.Level = slog.LevelWarn
	default:
		handlerOptions.Level = slog.LevelError
	}

	loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)
	slog.SetDefault(slog.New(loggerHandler))
}
