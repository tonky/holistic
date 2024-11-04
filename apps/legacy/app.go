package legacy

import (
	"context"
	"tonky/holistic/domain/food"
)

func (a App) ReadOrder(ctx context.Context, in food.OrderID) (food.Order, error) {
	return a.Deps.OrdererRepo.ReadOrderByID(ctx, in)
}

func (a App) CreateOrder(ctx context.Context, in NewOrder) (food.Order, error) {
	a.Deps.Logger.Info("legacy.App.CreateOrder", in)

	newOrder, err := a.Deps.OrdererRepo.SaveOrder(ctx, in)
	if err != nil {
		return food.Order{}, err
	}

	return newOrder, nil
}

func (a App) UpdateOrder(ctx context.Context, in UpdateOrder) (food.Order, error) {
	a.Deps.Logger.Info("legacy.App.UpdateOrder", in)

	updatedOrder, err := a.Deps.OrdererRepo.UpdateOrder(ctx, in)
	if err != nil {
		return food.Order{}, err
	}

	accOrder, errRO := a.Clients.AccountingClient.ReadOrder(ctx, updatedOrder.ID)
	if errRO != nil {
		a.Deps.Logger.Error("legacy.App.UpdateOrder accounting.Client.ReadOrder error", err, updatedOrder.ID)
	}

	a.Deps.Logger.Debug("legacy.App.UpdateOrder accounting.Client.ReadOrder ok", accOrder)

	if updatedOrder.IsFinal {
		if err := a.Deps.FoodOrderCreatedProducer.ProduceFoodOrderCreated(ctx, updatedOrder); err != nil {
			return food.Order{}, err
		}

		a.Deps.Logger.Debug("legacy.App.UpdateOrder published final order ok")
	}

	return updatedOrder, nil
}

type NewOrder struct {
	Content string
}

type UpdateOrder struct {
	ID      food.OrderID
	Content string
	IsFinal bool
}

func NewDeps(conf Config) Deps {
	return Deps{
		Config: conf,
	}
}
