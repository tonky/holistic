// AUTOGENERATED! DO NOT EDIT.
package kafkaConsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/kafka"

	"tonky/holistic/domain/food"
)

// compile-time check to make sure app-level interface is implemented
var _ IFoodOrderUpdated = new(FoodOrderUpdated) 

type IFoodOrderUpdated interface {
	Run(context.Context, func(context.Context, food.Order) error) chan error
}

type FoodOrderUpdated struct {
	logger logger.Slog
	client IConsumer
}

func NewFoodOrderUpdatedConsumer(logger logger.Slog, config kafka.Config) (*FoodOrderUpdated, error) {
	logger.Info(">> NewFoodOrderUpdatedConsumer()", "food.order.updated", config.GroupID)

	client := NewConsumer(config, "food.order.updated")

	return &FoodOrderUpdated{
		logger: logger,
		client: client,
	}, nil
}

func (c FoodOrderUpdated) Run(ctx context.Context, processor func(context.Context, food.Order) error) chan error {
	c.logger.Info(">> FoodOrderUpdated.Run()", c.client.Topic())

	res := make(chan error)
	models, errors := ConsumeFoodOrderUpdated(ctx, c.client)

	go func() {
		for {
			select {
			case model := <-models:
				c.logger.Info("kafkaConsumer.FoodOrderUpdated got model", model)

				if err := processor(ctx, model); err != nil {
					res <- err
				}
			case err := <-errors:
				res <- err
			case <-ctx.Done():
				return
			}
		}
	}()

	return res
}

func ConsumeFoodOrderUpdated(ctx context.Context, client IConsumer) (chan food.Order, chan error) {
	fmt.Println(">> ConsumeFoodOrderUpdated", client.Topic())

	models := make(chan food.Order)
	errors := make(chan error)

	kafkaMessages, kafkaErrors := client.Consume(context.Background())

	go func() {
		for {
			select {
			case err := <-kafkaErrors:
				errors <- err
			case <-ctx.Done():
				close(models)
				return
			case msg := <-kafkaMessages:
				var model food.Order
				if err := json.Unmarshal(msg.Value, &model); err != nil {
					errors <- err
					continue
				}
				models <- model
			}
		}
	}()

	return models, errors
}
