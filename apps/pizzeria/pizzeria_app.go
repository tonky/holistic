package pizzeria

import (
	"context"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra/kafkaProducer"

	"github.com/samber/do/v2"
)

func (app App) ReadOrder(ctx context.Context, in food.OrderID) (food.Order, error) {
	app.logger.Info("App.ReadOrder", in)

	or := do.MustInvokeAs[OrdererRepository](app.deps)

	return or.ReadOrderByID(ctx, in)
}

func (app App) CreateOrder(ctx context.Context, in NewOrder) (food.Order, error) {
	app.logger.Info("Pizzeria.CreateOrder", in)

	or := do.MustInvokeAs[OrdererRepository](app.deps)
	pr := do.MustInvokeAs[kafkaProducer.IFoodOrderCreated](app.deps)

	newOrder, err := or.SaveOrder(ctx, in)
	if err != nil {
		return food.Order{}, err
	}

	if err := pr.ProduceFoodOrderCreated(ctx, newOrder); err != nil {
		return food.Order{}, err
	}

	return newOrder, nil
}

func (app App) UpdateOrder(ctx context.Context, in UpdateOrder) (food.Order, error) {
	app.logger.Info("Pizzeria.UpdateOrder", in)

	// or := do.MustInvokeAs[OrdererRepository](app.deps)
	// pr := do.MustInvokeAs[kafkaProducer.IFoodOrderUpdated](app.deps)

	updatedOrder, err := app.ordererRepo.UpdateOrder(ctx, in)
	if err != nil {
		return food.Order{}, err
	}

	if err := app.FoodOrderUpdatedProducer.ProduceFoodOrderUpdated(ctx, updatedOrder); err != nil {
		return food.Order{}, err
	}

	return updatedOrder, nil
}

type NewOrder struct {
	Content string
}

type UpdateOrder struct {
	OrderID food.OrderID
	Content string
	IsFinal bool
}
