package main

import (
	app "tonky/holistic/apps/accounting"
	"tonky/holistic/infra/kafkaConsumer"
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

	consumer, err := kafkaConsumer.NewFoodOrderCreatedConsumer(logger, config.App.Kafka)
	if err != nil {
		panic(err)
	}

	injector := do.New()
	do.ProvideValue(injector, &config)
	do.ProvideValue(injector, &logger)

	// provide infra dependencies
	do.ProvideValue(injector, consumer)
	do.Provide(injector, app.NewOrdersMemoryRepository)

	svc, err := svc.New(injector)
	if err != nil {
		panic(err)
	}

	svc.Start()
}
