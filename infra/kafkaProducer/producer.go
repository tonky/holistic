package kafkaProducer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
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

	conn, err := kafka.DialLeader(context.Background(), "tcp", p.config.Brokers[0], p.topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: data},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	fmt.Println("Producer.Produce", p.topic, len(data), "done")
	return nil
}

func (p Producer) ProduceBatch(ctx context.Context, data [][]byte) error {
	fmt.Println("Producer.ProduceBatch", p.topic, len(data))

	return nil
}
