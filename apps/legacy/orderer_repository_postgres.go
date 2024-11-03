package legacy

import (
	"context"
	"tonky/holistic/domain/food"
)

func (p PostgresOrderer) ReadOrderByID(ctx context.Context, id food.OrderID) (food.Order, error) {
	return food.Order{}, nil
}

func (p PostgresOrderer) SaveOrder(ctx context.Context, in NewOrder) (food.Order, error) {
	return food.Order{}, nil
}

func (p PostgresOrderer) UpdateOrder(ctx context.Context, in UpdateOrder) (food.Order, error) {
	return food.Order{}, nil
}
