package main

import (
	app "tonky/holistic/apps/pizzeria"
	"tonky/holistic/infra/logger"
	svc "tonky/holistic/services/pizzeria"

	"github.com/samber/do/v2"
)

func main() {
	config, err := svc.NewEnvConfig()
	if err != nil {
		panic(err)
	}

	logger := logger.Slog{}

	por, err := app.NewPostgresOrdererRepository(logger, config.PostgresOrderer)
	if err != nil {
		panic(err)
	}

	kpfo, err := app.NewKafkaFoodOrderProducer(logger, config.KafkaProducerDefault)
	if err != nil {
		panic(err)
	}

	injector := do.New()
	do.ProvideValue(injector, &config)
	do.ProvideValue(injector, &logger)

	// provide infra dependencies
	do.ProvideValue(injector, kpfo)
	do.ProvideValue(injector, por)

	svc, err := svc.NewPizzeria(injector)
	if err != nil {
		panic(err)
	}

	svc.Start()
}
