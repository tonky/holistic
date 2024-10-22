package pizzeria

import (
	"context"
	"tonky/holistic/domain/food"
)

func (r PostgresOrderer) ReadOrderByID(ctx context.Context, req food.OrderID) (food.Order, error) {
	return food.Order{}, nil
}

func (r PostgresOrderer) SaveOrder(ctx context.Context, req NewOrder) (food.Order, error) {
	r.logger.Info("PostgresOrderer | CreateOrder", "pg conn", r.client, req)
	return food.Order{}, nil
}

type NewOrder struct{}
