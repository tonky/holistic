package pizzeria

import (
	"context"
	"fmt"
	"tonky/holistic/domain/food"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type MemoryOrderer struct {
	orders map[food.OrderID]food.Order
}

func (r MemoryOrderer) ReadOrderByID(ctx context.Context, req food.OrderID) (food.Order, error) {
	o, ok := r.orders[req]

	if ok {
		return o, nil
	}

	return food.Order{}, fmt.Errorf("order %s not found", req)
}

func (r MemoryOrderer) SaveOrder(ctx context.Context, req NewOrder) (food.Order, error) {
	order := food.Order{ID: food.OrderID{ID: uuid.New()}, Content: req.Content}

	r.orders[order.ID] = order

	return order, nil
}

func NewMemoryOrdererRepository(deps do.Injector) (*MemoryOrderer, error) {
	return &MemoryOrderer{orders: map[food.OrderID]food.Order{}}, nil
}
