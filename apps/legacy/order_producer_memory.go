package legacy

import (
	"context"
	"fmt"
	"tonky/holistic/domain/food"

	"github.com/google/uuid"
)

type MemoryOrderProducer struct {
	Orders map[food.OrderID]food.Order
}

func (r *MemoryOrderProducer) ProduceFoodOrderCreated(ctx context.Context, in food.Order) error {
	fmt.Println("MemoryOrderProducerRepository.ProduceNewOrder", in)

	order := food.Order{ID: food.OrderID{ID: uuid.New()}, Content: in.Content}

	r.Orders[order.ID] = order

	return nil
}

func (r *MemoryOrderProducer) ProduceFoodOrderCreatedBatch(ctx context.Context, in []food.Order) error {
	fmt.Println("MemoryOrderProducerRepository.ProduceFoodOrderBatch", in)
	return nil
}

func NewMemoryOrderProducer() *MemoryOrderProducer {
	store := map[food.OrderID]food.Order{}

	return &MemoryOrderProducer{Orders: store}
}
