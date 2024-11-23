// AUTOGENERATED! DO NOT EDIT.
package kafkaConsumer

import (
	"context"
	"encoding/json"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/kafkaConsumer"

	{% for i in topic.Obj.AbsImports(ctx) %}
	"{{ i }}"
	{% end %}
)

// compile-time check to make sure app-level interface is implemented
var _ {{ topic.InterfaceName() }} = new({{ topic.StructName() }}Consumer) 

type {{ topic.InterfaceName() }} interface {
	Run(context.Context, func(context.Context, {{ topic.ModelName() }}) error) chan error
}

type {{ topic.StructName() }}Consumer struct {
	logger logger.ILogger
	client kafkaConsumer.IConsumer
}

func New{{ topic.StructName() }}Consumer(l logger.ILogger, config kafka.Config) (*{{ topic.StructName() }}Consumer, error) {
	l = l.With("kafkaConsumer", "{{ topic.StructName() }}Consumer", "topic", "{{ topic.TopicName }}", "groupID", config.GroupID)
	
	l.Info("New consumer")

	client := kafkaConsumer.NewConsumer(config, "{{ topic.TopicName }}")

	return &{{ topic.StructName() }}Consumer {
		logger: l,
		client: client,
	}, nil
}

func (c {{ topic.StructName() }}Consumer) Run(ctx context.Context, processor func(context.Context, {{ topic.ModelName() }}) error) chan error {
	c.logger.Info("{{ topic.StructName() }}.Run()")

	res := make(chan error)
	models, errors := Consume{{ cap(topic.Name) }}(ctx, c.client)

	go func() {
		for {
			select {
			case model := <-models:
				c.logger.Info("got model", model)

				if err := processor(ctx, model); err != nil {
					res <- err
				}
			case err := <-errors:
				res <- err
			case <-ctx.Done():
				return
			}
		}
	}()

	return res
}

func Consume{{ cap(topic.Name) }}(ctx context.Context, client kafkaConsumer.IConsumer) (chan {{ topic.ModelName() }}, chan error) {
	client.Logger().Info("Consume{{ cap(topic.Name) }}", "topic", client.Topic())

	models := make(chan {{ topic.ModelName() }})
	errors := make(chan error)

	kafkaMessages, kafkaErrors := client.Consume(context.Background())

	go func() {
		for {
			select {
			case err := <-kafkaErrors:
				errors <- err
			case <-ctx.Done():
				close(models)
				return
			case msg := <-kafkaMessages:
				var model {{ topic.ModelName() }}
				if err := json.Unmarshal(msg.Value, &model); err != nil {
					errors <- err
					continue
				}
				models <- model
			}
		}
	}()

	return models, errors
}
