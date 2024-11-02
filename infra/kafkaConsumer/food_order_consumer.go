package kafkaConsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"tonky/holistic/domain/food"
)

func ConsumeFoodOrder(consumer IConsumer) (chan food.Order, chan error) {
	resModels := make(chan food.Order)
	resErrors := make(chan error)

	if consumer.Topic() != "pizzeria.order" {
		resErrors <- fmt.Errorf("expected pizzeria.food topic, got '%s'", consumer.Topic())
	}

	kafkaMessages, kafkaErrors := consumer.Consume(context.Background())

	go func() {
		consumer.Logger().Info("infra.kafkaConsumer.ConsumeFoodOrder | started goroutine")
		for {
			// consumer.Logger().Info("infra.kafkaConsumer.ConsumeFoodOrder | for loop...")
			select {
			case err := <-kafkaErrors:
				consumer.Logger().Error("infra.kafkaConsumer.ConsumeFoodOrder got error from kafka", err)
				resErrors <- fmt.Errorf("failed to consume message: %w", err)
			case m := <-kafkaMessages:
				var order food.Order
				err := json.Unmarshal(m.Value, &order)

				if err != nil {
					// consumer.Logger().Error("infra.kafkaConsumer.ConsumeFoodOrder | failed to unmarshal message", err, m.Value)
					resErrors <- fmt.Errorf("failed to unmarshal '%s' as food.Order from pizzeria.food topic: %w", string(m.Value), err)
					continue
				}

				consumer.Logger().Debug("infra.kafkaConsumer.ConsumeFoodOrder | got order from kafka, returning model", order)
				resModels <- order
			}
		}
	}()

	return resModels, resErrors
}
