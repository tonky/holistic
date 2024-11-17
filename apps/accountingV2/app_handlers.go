package accountingV2

import (
	"context"
	"tonky/holistic/domain/food"
	"tonky/holistic/domain/foodStore"
)

func (a App) GetOrderByID(ctx context.Context, in foodStore.OrderID) (food.Order, error) {
	return food.Order{}, nil
}

func (a App) CreateOrder(ctx context.Context, in NewFoodOrder) (food.Order, error) {
	return food.Order{}, nil
}
