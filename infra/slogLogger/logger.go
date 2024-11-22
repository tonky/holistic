package slogLogger

import (
	"fmt"
	"os"
	"strings"
	"time"
	"tonky/holistic/infra/logger"
)

var logLevel = strings.ToLower(os.Getenv("LOG_LEVEL"))

func Default() Slog {
	if logLevel == "" {
		logLevel = "info"
	}

	return Slog{}
}

func NewFromConfig(c logger.Config) Slog {
	logLevel = strings.ToLower(c.Level)

	return Slog{}
}

type Slog struct{}

func (sl Slog) Debug(msg string, fields ...interface{}) {
	if logLevel == "info" || logLevel == "warn" || logLevel == "error" {
		return
	}

	fmt.Printf(">> %s SlogLogger.Debug | %s | %s \n", time.Now().Format("2006-01-02 15:04:05.999"), msg, fields)

}

func (sl Slog) Info(msg string, fields ...interface{}) {
	if logLevel == "warn" || logLevel == "error" {
		return
	}

	fmt.Printf(">> %s SlogLogger.Info | %s | %s \n", time.Now().Format("2006-01-02 15:04:05.999"), msg, fields)
}

func (sl Slog) Warn(msg string, fields ...interface{}) {
	if logLevel == "error" {
		return
	}

	fmt.Printf(">> %s SlogLogger.Warn | %s | %s \n", time.Now().Format("2006-01-02 15:04:05.999"), msg, fields)
}

func (sl Slog) Error(msg string, fields ...interface{}) {
	fmt.Printf(">> %s SlogLogger.Error | %s | %s \n", time.Now().Format("2006-01-02 15:04:05.999"), msg, fields)
}
