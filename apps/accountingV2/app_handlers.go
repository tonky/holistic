package accountingV2

import (
	"context"
	"tonky/holistic/domain/accountingV2"
	"tonky/holistic/domain/foodStore"
)

func (a App) GetOrderByID(ctx context.Context, in foodStore.OrderID) (foodStore.Order, error) {
	return foodStore.Order{}, nil
}

func (a App) CreateOrder(ctx context.Context, in accountingV2.NewFoodOrder) (foodStore.Order, error) {
	return foodStore.Order{}, nil
}
