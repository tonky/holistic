package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/samber/do/v2"
)

var logLevel = strings.ToLower(os.Getenv("LOG_LEVEL"))

func NewSlogLogger(i do.Injector) (*Slog, error) {
	return &Slog{}, nil
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
