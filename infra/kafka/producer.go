package kafka

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/slogLogger"

	"github.com/segmentio/kafka-go"
)

type IProducer interface {
	Produce(context.Context, []byte) error
	ProduceBatch(context.Context, [][]byte) error
}

type Producer struct {
	config Config
	topic  string
	logger logger.ILogger
}

func NewProducer(config Config, topic string) Producer {
	return Producer{
		config: config,
		topic:  topic,
		logger: slogLogger.Default(),
	}
}

func (p Producer) Produce(ctx context.Context, data []byte) error {
	slog.Info("Producer.Produce", slog.String("topic", p.topic), slog.Int("data len", len(data)))

	w := &kafka.Writer{
		Addr:                   kafka.TCP(p.config.Brokers...),
		Topic:                  p.topic,
		AllowAutoTopicCreation: true,
		BatchSize:              1,
	}

	err := w.WriteMessages(ctx, kafka.Message{Value: data})
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	/*
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
	*/

	p.logger.Info("Producer.Produce OK", "topic", p.topic, "data len", len(data))

	return nil
}

func (p Producer) ProduceBatch(ctx context.Context, data [][]byte) error {
	fmt.Println("Producer.ProduceBatch", p.topic, len(data))

	return nil
}
