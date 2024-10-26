package pizzeria

import (
	"context"
	"fmt"
	"tonky/holistic/domain/food"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type MemoryOrderProducerRepository struct {
	orders map[food.OrderID]food.Order
}

func (r MemoryOrderProducerRepository) ProduceFoodOrder(ctx context.Context, in food.Order) error {
	fmt.Println("MemoryOrderProducerRepository.ProduceNewOrder", in)

	order := food.Order{ID: food.OrderID{ID: uuid.New()}, Content: in.Content}

	r.orders[order.ID] = order

	return nil
}

func (r MemoryOrderProducerRepository) ProduceFoodOrderBatch(ctx context.Context, in []food.Order) error {
	fmt.Println("MemoryOrderProducerRepository.ProduceFoodOrderBatch", in)
	return nil
}

func NewMemoryOrderProducerRepository(deps do.Injector) (*MemoryOrderProducerRepository, error) {
	return &MemoryOrderProducerRepository{orders: map[food.OrderID]food.Order{}}, nil
}
