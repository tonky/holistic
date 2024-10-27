package logger

import (
	"fmt"
	"time"

	"github.com/samber/do/v2"
)

func NewSlogLogger(i do.Injector) (*Slog, error) {
	return &Slog{}, nil
}

type Slog struct{}

func (sl Slog) Info(msg string, fields ...interface{}) {
	fmt.Printf(">> %s SlogLogger.Info | %s | %s \n", time.Now().Format("2006-01-02 15:04:05.999"), msg, fields)
}

func (sl Slog) Debug(msg string, fields ...interface{}) {}
