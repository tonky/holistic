package main

import (
	app "tonky/holistic/apps/accounting"
	"tonky/holistic/infra/logger"
	svc "tonky/holistic/services/accounting"

	"github.com/samber/do/v2"
)

func main() {
	config, err := svc.NewEnvConfig()
	if err != nil {
		panic(err)
	}

	logger := logger.Slog{}

	consumer, err := app.NewKafkaFoodOrderConsumer(logger, config.Kafka)
	if err != nil {
		panic(err)
	}

	injector := do.New()
	do.ProvideValue(injector, &config)
	do.ProvideValue(injector, &logger)

	// provide infra dependencies
	do.ProvideValue(injector, consumer)

	svc, err := svc.NewAccounting(injector)
	if err != nil {
		panic(err)
	}

	svc.Start()
}
