package otelinit

import (
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/slogLogger"
	"tonky/holistic/infra/tele"
)

func NewFromConfig(conf Config, serviceName, packagePath string) (*tele.Otel, error) {
	metricer, _, err := InitMetrics(conf, "accountingV2")
	if err != nil {
		return nil, err
	}

	tracer, _, err := InitTracing(conf, "accountingV2", "path/to/accountingV2")
	if err != nil {
		return nil, err
	}

	return &tele.Otel{
		Logger:  slogLogger.Default(),
		Metrics: metricer,
		Tracer:  tracer,
	}, nil
}

func NewLogger() logger.ILogger {
	return slogLogger.Default()
}
