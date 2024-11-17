// AUTOGENERATED! DO NOT EDIT.
package kafkaProducer

import (
	"context"
	"encoding/json"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/kafkaProducer"

	{% for i in topic.Obj.AbsImports(ctx) %}
	"{{ i }}"
	{% end %}
)

// compile-time check to make sure app-level interface is implemented
var _ {{ topic.InterfaceName() }} = new({{ topic.StructName() }}Producer) 

type {{ topic.InterfaceName() }} interface {
	Produce{{ cap(topic.Name) }}(context.Context, {{ topic.ModelName() }}) error
	Produce{{ cap(topic.Name) }}Batch(context.Context, []{{ topic.ModelName() }}) error
}

type {{ topic.StructName() }}Producer struct {
	logger logger.Slog
	client kafkaProducer.IProducer
}

func New{{ topic.StructName() }}Producer(logger logger.Slog, config kafka.Config) (*{{ topic.StructName() }}Producer, error) {
	client := kafkaProducer.NewProducer(config, "{{ topic.TopicName }}")

	return &{{ topic.StructName() }}Producer {
		logger: logger,
		client: client,
	}, nil
}

func (r {{ topic.StructName() }}Producer) Produce{{ cap(topic.Name) }}(ctx context.Context, in {{ topic.ModelName() }}) error {
	r.logger.Info("{{ topic.StructName() }}.Produce{{ cap(topic.Name) }}", in)

	inBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return r.client.Produce(ctx, inBytes)
}

func (r {{ topic.StructName() }}Producer) Produce{{ cap(topic.Name) }}Batch(ctx context.Context, ins []{{ topic.ModelName() }}) error {
	r.logger.Info("{{ topic.StructName() }}.Produce{{ cap(topic.Name) }}Batch", ins)

	var data [][]byte

	for _, in  := range ins {
		inBytes, err := json.Marshal(in)
		if err != nil {
			return err
		}
	
		data = append(data, inBytes)
	}

	return r.client.ProduceBatch(ctx, data)
}
