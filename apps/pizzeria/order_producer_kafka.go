package pizzeria

import (
	"context"
	"tonky/holistic/domain/food"
)

func (r KafkaFoodOrderProducer) ProduceOrder(ctx context.Context, in food.Order) error {
	r.logger.Info("KafkaOrderProducer.ProduceOrder", in)

	return r.client.Produce(ctx, []byte(in.Content))
}

func (r KafkaFoodOrderProducer) ProduceOrderBatch(ctx context.Context, in []food.Order) error {
	r.logger.Info("KafkaOrderProducer.ProduceOrderBatch", in)

	var data [][]byte
	for _, order := range in {
		data = append(data, []byte(order.Content))
	}

	return r.client.ProduceBatch(ctx, data)
}
