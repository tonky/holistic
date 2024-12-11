package tele

import (
	"fmt"
	"tonky/holistic/infra/logger"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Otel struct {
	Logger  logger.ILogger
	Metrics metric.Meter
	Tracer  trace.Tracer
}

func NewFromEnv() Otel {
	fmt.Println("tele.NewFromEnv()")
	return Otel{}
}
