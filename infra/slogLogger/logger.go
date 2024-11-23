package slogLogger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"tonky/holistic/infra/logger"
)

func Default() Slog {
	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))

	if logLevel == "" {
		fmt.Printf("infra.slogLogger.Default(): 'LOG_LEVEL' is empty, setting to 'info'\n")

		logLevel = "info"
	}

	return NewFromConfig(logger.Config{Level: logLevel})
}

func NewFromConfig(c logger.Config) Slog {
	slogLevel := getLevel(c.Level)

	lvl := new(slog.LevelVar)

	lvl.Set(slogLevel)

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))

	return Slog{logger: *logger}
}

func getLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type Slog struct {
	logger slog.Logger
}

func (sl Slog) Debug(msg string, fields ...interface{}) {
	sl.logger.Debug(msg, fields...)

}

func (sl Slog) Info(msg string, fields ...interface{}) {
	sl.logger.Info(msg, fields...)
}

func (sl Slog) Warn(msg string, fields ...interface{}) {
	sl.logger.Warn(msg, fields...)
}

func (sl Slog) Error(msg string, fields ...interface{}) {
	sl.logger.Error(msg, fields...)
}

func (sl Slog) With(key string, value any) logger.ILogger {
	sl.logger = *sl.logger.With(key, value)

	return sl
}
