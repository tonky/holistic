package pizzeria

import (
	"context"
	"tonky/holistic/domain/food"
)

func (app App) ReadOrder(ctx context.Context, in food.OrderID) (food.Order, error) {
	app.logger.Info("App.ReadOrder", in)

	return app.ordererRepo.ReadOrderByID(ctx, in)
}

func (app App) CreateOrder(ctx context.Context, in NewOrder) (food.Order, error) {
	app.logger.Info("App.CreateOrder", in)

	newOrder, err := app.ordererRepo.SaveOrder(ctx, in)
	if err != nil {
		return food.Order{}, err
	}

	if err := app.foodOrderProducer.ProduceFoodOrder(ctx, newOrder); err != nil {
		return food.Order{}, err
	}

	return newOrder, nil
}
