package kafka

import (
	"context"

	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/slogLogger"

	"github.com/segmentio/kafka-go"
)

type IConsumer interface {
	Consume(context.Context) (chan kafka.Message, chan error)
	Topic() string
	Logger() logger.ILogger
}

type Consumer struct {
	config Config
	topic  string
	logger logger.ILogger
}

func NewConsumer(config Config, topic string) Consumer {
	return Consumer{
		config: config,
		topic:  topic,
		logger: slogLogger.Default(),
	}
}

func (c Consumer) Topic() string {
	return c.topic
}

func (c Consumer) Logger() logger.ILogger {
	return c.logger
}

func (c Consumer) Consume(ctx context.Context) (chan kafka.Message, chan error) {
	c.logger.Info("infra.kafkaConsumer.Consumer.Consume starting for", "topic", c.topic, "conf", c.config)

	resMessages := make(chan kafka.Message)
	resErrors := make(chan error)

	reader := getKafkaReader(c.config.Brokers, c.topic, c.config.GroupID)
	reader.SetOffset(kafka.LastOffset)

	go func() {
		for {
			c.logger.Debug("infra.kafkaConsumer.Consumer.Consume | inner loop from", "topic", c.topic)
			m, err := reader.ReadMessage(ctx)

			c.logger.Debug("infra.kafkaConsumer.Consumer.Consume | inner loop got message from", "topic", c.topic)

			if err != nil {
				c.logger.Error("infra.kafkaConsumer.Consumer.Consume | consumer error, exiting: %v\n", "error", err)

				resErrors <- err

				if err := reader.Close(); err != nil {
					c.logger.Error("infra.kafkaConsumer.Consumer.Consume | failed to close reader: %v\n", "error", err)
				}

				return
			}

			c.logger.Debug("infra.kafkaConsumer.Consumer.Consume | consumed ", "topic", m.Topic, "partition", m.Partition, "offset", m.Offset, "key", string(m.Key))

			resMessages <- m
		}
	}()

	return resMessages, resErrors
}

func getKafkaReader(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		// GroupID:     groupID,
		Topic:       topic,
		GroupTopics: []string{topic},
		StartOffset: kafka.FirstOffset,
		// ReadBackoffMin:        10 * time.Millisecond,
		// WatchPartitionChanges: true,
		// MinBytes: 1, // 10KB
		// MaxBytes: 10e6, // 10MB
	})
}
