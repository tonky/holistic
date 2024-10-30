package pizzeria

import (
	"context"
	"encoding/json"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

func (r KafkaFoodOrderProducer) ProduceFoodOrder(ctx context.Context, in food.Order) error {
	r.logger.Info("KafkaFoodOrderProducer.ProduceOrder", in)

	orderData, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return r.client.Produce(ctx, orderData)
}

func (r KafkaFoodOrderProducer) ProduceFoodOrderBatch(ctx context.Context, in []food.Order) error {
	r.logger.Info("KafkaFoodOrderProducer.ProduceOrderBatch", in)

	var data [][]byte
	for _, order := range in {
		data = append(data, []byte(order.Content))
	}

	return r.client.ProduceBatch(ctx, data)
}

func NewDOKafkaFoodOrderProducer(deps do.Injector) (*KafkaFoodOrderProducer, error) {
	config := do.MustInvoke[*kafka.Config](deps)
	logger := do.MustInvoke[*logger.Slog](deps)

	return NewKafkaFoodOrderProducer(*logger, *config)
}
