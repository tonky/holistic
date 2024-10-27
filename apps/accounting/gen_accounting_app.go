// AUTOGENERATED! DO NOT EDIT.

package accounting

import (
	"context"
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

type App struct {
	deps       do.Injector
	logger     *logger.Slog

    foodOrderConsumer FoodOrderConsumer
}

func NewApp(deps do.Injector) (*App, error) {
	ctx := context.Background()

	app := App{
		deps:       deps,
		logger:     do.MustInvoke[*logger.Slog](deps),
        foodOrderConsumer: do.MustInvokeAs[FoodOrderConsumer](deps),
	}

	go app.foodOrderConsumer.Run(ctx, app.FoodOrderProcessor)
	return &app, nil
}
