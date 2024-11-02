package kafkaConsumer

import (
	"context"
	"strings"

	kafkaInfra "tonky/holistic/infra/kafka"
	"tonky/holistic/infra/logger"

	"github.com/segmentio/kafka-go"
)

type IConsumer interface {
	Consume(context.Context) (chan kafka.Message, chan error)
	Topic() string
	Logger() logger.Slog
}

type Consumer struct {
	config kafkaInfra.Config
	topic  string
	logger logger.Slog
}

func NewConsumer(config kafkaInfra.Config, topic string) Consumer {
	return Consumer{
		config: config,
		topic:  topic,
		logger: logger.Slog{},
	}
}

func (c Consumer) Topic() string {
	return c.topic
}

func (c Consumer) Logger() logger.Slog {
	return c.logger
}

func (c Consumer) Consume(ctx context.Context) (chan kafka.Message, chan error) {
	c.logger.Info("infra.kafkaConsumer.Consumer.Consume starting for", c.topic)

	resMessages := make(chan kafka.Message)
	resErrors := make(chan error)

	reader := getKafkaReader(c.config.Brokers[0], c.topic, "accounting")

	c.logger.Debug("infra.kafkaConsumer.Consumer.Consume | start consuming from %+v\n", reader.Config())

	go func() {
		for {
			c.logger.Debug("infra.kafkaConsumer.Consumer.Consume | reading messages...")
			m, err := reader.ReadMessage(ctx)

			c.logger.Debug("infra.kafkaConsumer.Consumer.Consume | got message!")

			if err != nil {
				c.logger.Error("infra.kafkaConsumer.Consumer.Consume | consumer error, exiting: %v\n", err)

				resErrors <- err

				if err := reader.Close(); err != nil {
					c.logger.Error("infra.kafkaConsumer.Consumer.Consume | failed to close reader: %v\n", err)
				}

				return
			}

			c.logger.Debug("infra.kafkaConsumer.Consumer.Consume | message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

			resMessages <- m
		}
	}()

	return resMessages, resErrors
}

func getKafkaReader(kafkaURL, topic, _ string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		// GroupID: groupID,
		Topic: topic,
		// MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}
