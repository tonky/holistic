package accounting

import (
	"context"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/kafkaConsumer"
	"tonky/holistic/infra/logger"

	"tonky/holistic/domain/food"
)

// compile-time check to make sure app-level interface is implemented
var _ FoodOrderConsumer = new(KafkaFoodOrderConsumer)

type FoodOrderConsumer interface {
	Run(context.Context, func(context.Context, food.Order) error) chan error
}

type KafkaFoodOrderConsumer struct {
	logger logger.Slog
	client kafkaConsumer.IConsumer
}

func NewKafkaFoodOrderConsumer(logger logger.Slog, config kafka.Config) (*KafkaFoodOrderConsumer, error) {
	client := kafkaConsumer.NewConsumer(config, "food.order.created")

	return &KafkaFoodOrderConsumer{
		logger: logger,
		client: client,
	}, nil
}

func (c KafkaFoodOrderConsumer) Run(ctx context.Context, processor func(context.Context, food.Order) error) chan error {
	res := make(chan error)
	orders, errors := kafkaConsumer.ConsumeFoodOrder(c.client)

	go func() {
		for {
			select {
			case order := <-orders:
				c.logger.Info("accounting.KafkaFoodOrderConsumer got order", order)

				if err := processor(ctx, order); err != nil {
					res <- err
				}
			case err := <-errors:
				res <- err
			}
		}
	}()

	return res
}
