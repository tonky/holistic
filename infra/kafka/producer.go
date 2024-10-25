package kafka

import "context"

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
