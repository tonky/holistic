package accounting

import (
	"context"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"

	"github.com/samber/do/v2"
)

var _ OrdererRepository = new(OrdersMemoryRepository)

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

func (a OrdersMemoryRepository) SaveFinishedOrder(ctx context.Context, in accounting.Order) (accounting.Order, error) {
	a.orders[in.ID] = in

	return in, nil
}
