package kafkaProducer

import (
	"context"
	"fmt"
	"log"
	"time"

	kafkaInfra "tonky/holistic/infra/kafka"
	"tonky/holistic/infra/logger"

	"github.com/segmentio/kafka-go"
)

type IProducer interface {
	Produce(context.Context, []byte) error
	ProduceBatch(context.Context, [][]byte) error
}

type Producer struct {
	config kafkaInfra.Config
	topic  string
	logger *logger.Slog
}

func NewProducer(config kafkaInfra.Config, topic string) Producer {
	return Producer{
		config: config,
		topic:  topic,
		logger: &logger.Slog{},
	}
}

func (p Producer) Produce(ctx context.Context, data []byte) error {
	p.logger.Info("Producer.Produce", p.topic, len(data))

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

	p.logger.Info("Producer.Produce", p.topic, len(data), "done")
	return nil
}

func (p Producer) ProduceBatch(ctx context.Context, data [][]byte) error {
	fmt.Println("Producer.ProduceBatch", p.topic, len(data))

	return nil
}
