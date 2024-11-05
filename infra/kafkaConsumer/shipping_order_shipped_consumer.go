// AUTOGENERATED! DO NOT EDIT.
package kafkaConsumer

import (
	"context"
	"encoding/json"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/kafka"

	"tonky/holistic/domain/shipping"
)

// compile-time check to make sure app-level interface is implemented
var _ IShippingOrderShipped = new(ShippingOrderShipped) 

type IShippingOrderShipped interface {
	Run(context.Context, func(context.Context, shipping.Order) error) chan error
}

type ShippingOrderShipped struct {
	logger logger.Slog
	client IConsumer
}

func NewShippingOrderShippedConsumer(logger logger.Slog, config kafka.Config) (*ShippingOrderShipped, error) {
	client := NewConsumer(config, "shipping.order.shipped")

	return &ShippingOrderShipped{
		logger: logger,
		client: client,
	}, nil
}

func (c ShippingOrderShipped) Run(ctx context.Context, processor func(context.Context, shipping.Order) error) chan error {
	res := make(chan error)
	models, errors := ConsumeShippingOrderShipped(ctx, c.client)

	go func() {
		for {
			select {
			case model := <-models:
				c.logger.Info("kafkaConsumer.ShippingOrderShipped got model", model)

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

func ConsumeShippingOrderShipped(ctx context.Context, client IConsumer) (chan shipping.Order, chan error) {
	models := make(chan shipping.Order)
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
				var model shipping.Order
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