package pizzeria

import (
	"context"
	"tonky/holistic/domain/food"
)

func (r KafkaFoodOrderProducer) ProduceFoodOrder(ctx context.Context, in food.Order) error {
	r.logger.Info("KafkaFoodOrderProducer.ProduceOrder", in)

	return r.client.Produce(ctx, []byte(in.ID.ID.String()+in.Content))
}

func (r KafkaFoodOrderProducer) ProduceFoodOrderBatch(ctx context.Context, in []food.Order) error {
	r.logger.Info("KafkaFoodOrderProducer.ProduceOrderBatch", in)

	var data [][]byte
	for _, order := range in {
		data = append(data, []byte(order.Content))
	}

	return r.client.ProduceBatch(ctx, data)
}
