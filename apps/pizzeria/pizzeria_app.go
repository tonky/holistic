package pizzeria

import (
	"context"
	"log/slog"
	"tonky/holistic/domain/food"
)

func (app App) ReadOrder(ctx context.Context, id food.OrderID) (food.Order, error) {
	app.logger.Info("ReadOrder pgxConn", slog.Any("pg conn", app.ordererRepo))

	res := food.Order{ID: id}
	return res, nil
}

func (app App) CreateOrder(ctx context.Context, in food.Order) (food.Order, error) {
	// var res food.Order
	res, err := app.ordererRepo.SaveOrder(ctx, NewOrder{})

	return res, err
}
