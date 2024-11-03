package legacy

import (
	"context"
	"fmt"
	"tonky/holistic/domain/food"

	"github.com/google/uuid"
)

type MemoryOrderProducer struct {
	orders map[food.OrderID]food.Order
}

func (r *MemoryOrderProducer) ProduceFoodOrder(ctx context.Context, in food.Order) error {
	fmt.Println("MemoryOrderProducerRepository.ProduceNewOrder", in)

	order := food.Order{ID: food.OrderID{ID: uuid.New()}, Content: in.Content}

	r.orders[order.ID] = order

	return nil
}

func (r *MemoryOrderProducer) ProduceFoodOrderBatch(ctx context.Context, in []food.Order) error {
	fmt.Println("MemoryOrderProducerRepository.ProduceFoodOrderBatch", in)
	return nil
}

func NewMemoryOrderProducer() (*MemoryOrderProducer, *map[food.OrderID]food.Order) {
	store := map[food.OrderID]food.Order{}

	return &MemoryOrderProducer{orders: store}, &store
}
