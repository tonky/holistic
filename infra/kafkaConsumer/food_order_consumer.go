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
		fmt.Println("infra.kafkaConsumer.ConsumeFoodOrder | started goroutine")
		for {
			select {
			case err := <-kafkaErrors:
				fmt.Println("infra.kafkaConsumer.ConsumeFoodOrder got error from kafka", err)
				resErrors <- fmt.Errorf("failed to consume message: %w", err)
			case m := <-kafkaMessages:
				var order food.Order
				err := json.Unmarshal(m.Value, &order)

				if err != nil {
					fmt.Println("infra.kafkaConsumer.ConsumeFoodOrder | failed to unmarshal message", err)
					resErrors <- fmt.Errorf("failed to unmarshal message order from pizzeria.food topic: %w", err)
					continue
				}

				fmt.Println("infra.kafkaConsumer.ConsumeFoodOrder | got order from kafka, returning model", order)
				resModels <- order
			}
		}
	}()

	return resModels, resErrors
}
