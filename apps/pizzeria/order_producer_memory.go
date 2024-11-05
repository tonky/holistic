package pizzeria

import (
	"context"
	"fmt"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra/kafkaProducer"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

var _ kafkaProducer.IFoodOrderCreated = new(MemoryOrderProducerRepository)

type MemoryOrderProducerRepository struct {
	orders map[food.OrderID]food.Order
}

func (r MemoryOrderProducerRepository) ProduceFoodOrderCreated(ctx context.Context, in food.Order) error {
	fmt.Println("MemoryOrderProducerRepository.ProduceNewOrder", in)

	order := food.Order{ID: food.OrderID{ID: uuid.New()}, Content: in.Content}

	r.orders[order.ID] = order

	return nil
}

func (r MemoryOrderProducerRepository) ProduceFoodOrderCreatedBatch(ctx context.Context, in []food.Order) error {
	fmt.Println("MemoryOrderProducerRepository.ProduceFoodOrderBatch", in)
	return nil
}

func NewMemoryOrderProducerRepository(deps do.Injector) (*MemoryOrderProducerRepository, error) {
	return &MemoryOrderProducerRepository{orders: map[food.OrderID]food.Order{}}, nil
}
