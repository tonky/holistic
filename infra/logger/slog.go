package logger

import (
	"fmt"

	"github.com/samber/do/v2"
)

func NewSlogLogger(i do.Injector) (*Slog, error) {
	return &Slog{}, nil
}

type Slog struct{}

func (sl Slog) Info(msg string, fields ...interface{}) {
	fmt.Println(">> SlogLogger.Info | ", msg, fields)
}

func (sl Slog) Debug(msg string, fields ...interface{}) {}
