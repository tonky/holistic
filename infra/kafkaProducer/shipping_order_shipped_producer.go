// AUTOGENERATED! DO NOT EDIT.
package kafkaProducer

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
	ProduceShippingOrderShipped(context.Context, shipping.Order) error
	ProduceShippingOrderShippedBatch(context.Context, []shipping.Order) error
}

type ShippingOrderShipped struct {
	logger logger.ILogger
	client IProducer
}

func NewShippingOrderShippedProducer(l logger.ILogger, config kafka.Config) (*ShippingOrderShipped, error) {
	client := NewProducer(config, "shipping.order.shipped")

	return &ShippingOrderShipped{
		logger: l,
		client: client,
	}, nil
}

func (r ShippingOrderShipped) ProduceShippingOrderShipped(ctx context.Context, in shipping.Order) error {
	r.logger.Info("ShippingOrderShipped.ProduceShippingOrderShipped", in)

	inBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return r.client.Produce(ctx, inBytes)
}

func (r ShippingOrderShipped) ProduceShippingOrderShippedBatch(ctx context.Context, ins []shipping.Order) error {
	r.logger.Info("ShippingOrderShipped.ProduceShippingOrderShippedBatch", ins)

	var data [][]byte

	for _, in  := range ins {
		inBytes, err := json.Marshal(in)
		if err != nil {
			return err
		}
	
		data = append(data, inBytes)
	}

	return r.client.ProduceBatch(ctx, data)
}
