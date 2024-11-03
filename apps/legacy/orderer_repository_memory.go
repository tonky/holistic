package legacy

import (
	"context"
	"fmt"
	"tonky/holistic/domain/food"
)

var _ OrdererRepository = new(MemoryOrderer)

type MemoryOrderer struct {
	orders map[food.OrderID]food.Order
}

func (p *MemoryOrderer) ReadOrderByID(ctx context.Context, id food.OrderID) (food.Order, error) {
	order, ok := p.orders[id]

	if !ok {
		return food.Order{}, fmt.Errorf("order not found")
	}

	return order, nil
}

func (p *MemoryOrderer) SaveOrder(ctx context.Context, in NewOrder) (food.Order, error) {
	fmt.Println("MemoryOrderer.SaveOrder", in)

	order := food.Order{
		ID:      food.RandomOrderID(),
		Content: in.Content,
	}

	p.orders[order.ID] = order

	return order, nil
}

func (p *MemoryOrderer) UpdateOrder(ctx context.Context, in UpdateOrder) (food.Order, error) {
	order, ok := p.orders[in.ID]

	if !ok {
		return food.Order{}, fmt.Errorf("order not found")
	}

	order.Content = in.Content
	order.IsFinal = in.IsFinal

	p.orders[order.ID] = order

	return order, nil
}

func NewMemoryOrdererRepository() *MemoryOrderer {
	return &MemoryOrderer{
		orders: make(map[food.OrderID]food.Order),
	}
}
