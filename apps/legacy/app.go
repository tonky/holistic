package legacy

import (
	"context"
	accClient "tonky/holistic/clients/accounting"
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

	accountingClient := accClient.NewAccounting(accClient.Config{Host: "localhost", Port: 1235})

	_, errRO := accountingClient.ReadOrder(ctx, updatedOrder.ID)
	if errRO != nil {
		a.Deps.Logger.Error("legacy.App.UpdateOrder accounting.Client.ReadOrder error", err, updatedOrder.ID)
	}

	if updatedOrder.IsFinal {
		if err := a.Deps.FoodOrderProducer.ProduceFoodOrder(ctx, updatedOrder); err != nil {
			return food.Order{}, err
		}
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
