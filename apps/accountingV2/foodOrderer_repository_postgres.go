package accountingV2

import (
	"context"
	"tonky/holistic/domain/accountingV2"
	"tonky/holistic/domain/foodStore"
)

func (r FoodOrderer) GetOrderByID(ctx context.Context, orderID foodStore.OrderID) (foodStore.Order, error) {
	return foodStore.Order{}, nil
}

func (r FoodOrderer) SaveNewOrder(ctx context.Context, newOrder accountingV2.NewFoodOrder) (foodStore.Order, error) {
	return foodStore.Order{}, nil
}
