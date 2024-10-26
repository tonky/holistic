package kafka

import (
	"context"
	"fmt"
)

type IProducer interface {
	Produce(context.Context, []byte) error
	ProduceBatch(context.Context, [][]byte) error
}

type Producer struct {
	config Config
	topic  string
}

func NewProducer(config Config, topic string) Producer {
	return Producer{
		config: config,
		topic:  topic,
	}
}

func (p Producer) Produce(ctx context.Context, data []byte) error {
	fmt.Println("Producer.Produce", p.topic, len(data))

	return nil
}

func (p Producer) ProduceBatch(ctx context.Context, data [][]byte) error {
	fmt.Println("Producer.ProduceBatch", p.topic, len(data))

	return nil
}
