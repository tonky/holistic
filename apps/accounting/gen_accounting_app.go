package accounting

import (
	"context"
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

type App struct {
	deps   do.Injector
	logger *logger.Slog

	FoodOrderUpdatedConsumer FoodOrderConsumer
}

func NewApp(deps do.Injector) (*App, error) {
	ctx := context.Background()

	app := App{
		deps:                     deps,
		logger:                   do.MustInvoke[*logger.Slog](deps),
		FoodOrderUpdatedConsumer: do.MustInvokeAs[FoodOrderConsumer](deps),
	}

	go func() {
		for err := range app.FoodOrderUpdatedConsumer.Run(ctx, app.FoodOrderUpdatedProcessor) {
			app.logger.Warn(err.Error())
		}
	}()

	return &app, nil
}
