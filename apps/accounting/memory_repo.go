package accounting

import (
	"context"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"

	"github.com/samber/do/v2"
)

var _ AccountOrdersRepoReader = new(OrdersMemoryRepository)

func NewOrdersMemoryRepository(do.Injector) (*OrdersMemoryRepository, error) {
	return &OrdersMemoryRepository{
		orders: make(map[food.OrderID]accounting.Order),
	}, nil
}

type OrdersMemoryRepository struct {
	orders map[food.OrderID]accounting.Order
}

func (a OrdersMemoryRepository) ReadOrderByFoodID(ctx context.Context, in food.OrderID) (accounting.Order, error) {
	return a.orders[in], nil
}

func (a OrdersMemoryRepository) SaveOrder(ctx context.Context, in NewOrder) (accounting.Order, error) {
	out := accounting.Order{
		ID:   in.Order.ID,
		Cost: in.Cost,
	}

	a.orders[in.Order.ID] = out

	return out, nil
}
