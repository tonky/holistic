package accountingV2

import (
	"context"
	"tonky/holistic/domain/foodStore"
)

func (a App) GetOrderByID(ctx context.Context, in foodStore.OrderID) (foodStore.Order, error) {
	return foodStore.Order{}, nil
}

func (a App) CreateOrder(ctx context.Context, in NewFoodOrder) (foodStore.Order, error) {
	return foodStore.Order{}, nil
}
