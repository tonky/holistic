// AUTOGENERATED! DO NOT EDIT.
package {{ service_name }}

import (
	"context"
	"encoding/json"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/kafkaProducer"

	{% for imp in kp.Imports() %}
	{{ imp.Alias}} "tonky/holistic/{{ imp.RelPath }}"
	{% end %}
)

// compile-time check to make sure app-level interface is implemented
var _ {{ kp.InterfaceName() }} = new({{ kp.StructName() }}) 

type {{ kp.InterfaceName() }} interface {
	Produce{{ cap(kp.Name) }}(context.Context, {{ kp.Model }}) error
	Produce{{ cap(kp.Name) }}Batch(context.Context, []{{ kp.Model }}) error
}

type {{ kp.StructName() }} struct {
	logger logger.Slog
	client kafkaProducer.IProducer
}

func New{{ kp.StructName() }}(logger logger.Slog, config kafka.Config) (*{{ kp.StructName() }}, error) {
	client := kafkaProducer.NewProducer(config, "{{ kp.Topic }}")

	return &{{ kp.StructName() }}{
		logger: logger,
		client: client,
	}, nil
}

func (r {{ kp.StructName() }}) Produce{{ cap(kp.Name) }}(ctx context.Context, in {{ kp.Model }}) error {
	r.logger.Info("{{ kp.StructName() }}.Produce{{ cap(kp.Name) }}", in)

	inBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return r.client.Produce(ctx, inBytes)
}

func (r {{ kp.StructName() }}) Produce{{ cap(kp.Name) }}Batch(ctx context.Context, ins []{{ kp.Model }}) error {
	r.logger.Info("{{ kp.StructName() }}.Produce{{ cap(kp.Name) }}Batch", ins)

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
