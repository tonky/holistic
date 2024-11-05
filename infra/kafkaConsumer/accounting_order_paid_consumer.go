// AUTOGENERATED! DO NOT EDIT.
package kafkaConsumer

import (
	"context"
	"encoding/json"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/kafka"

	"tonky/holistic/domain/accounting"
)

// compile-time check to make sure app-level interface is implemented
var _ IAccountingOrderPaid = new(AccountingOrderPaid) 

type IAccountingOrderPaid interface {
	Run(context.Context, func(context.Context, accounting.Order) error) chan error
}

type AccountingOrderPaid struct {
	logger logger.Slog
	client IConsumer
}

func NewAccountingOrderPaidConsumer(logger logger.Slog, config kafka.Config) (*AccountingOrderPaid, error) {
	client := NewConsumer(config, "accounting.order.paid")

	return &AccountingOrderPaid{
		logger: logger,
		client: client,
	}, nil
}

func (c AccountingOrderPaid) Run(ctx context.Context, processor func(context.Context, accounting.Order) error) chan error {
	res := make(chan error)
	models, errors := ConsumeAccountingOrderPaid(ctx, c.client)

	go func() {
		for {
			select {
			case model := <-models:
				c.logger.Info("kafkaConsumer.AccountingOrderPaid got model", model)

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

func ConsumeAccountingOrderPaid(ctx context.Context, client IConsumer) (chan accounting.Order, chan error) {
	models := make(chan accounting.Order)
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
				var model accounting.Order
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
