package logger

import (
	"fmt"

	"github.com/samber/do/v2"
)

func NewSlogLogger(i do.Injector) (*SlogLogger, error) {
	return &SlogLogger{}, nil
}

type SlogLogger struct{}

func (sl SlogLogger) Info(msg string, fields ...interface{}) {
	fmt.Println(">> SlogLogger.Info | ", msg, fields)
}

func (sl SlogLogger) Debug(msg string, fields ...interface{}) {}
