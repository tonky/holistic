// AUTOGENERATED! DO NOT EDIT.
package {{ service.Name }}

import (
	"context"
	"{{ modulePath }}/infra/logger"
	"{{ modulePath }}/infra/postgres"
    {% for ci in service.Postgres.AbsImports(ctx) %}
	"{{ ci }}"
    {% end %}
)

var _ {{ repo.InterfaceName() }} = new({{ repo.StructName() }})

type {{ repo.InterfaceName() }} interface {
    {% for e in repo.Endpoints %}
    {{ e.Name }}(context.Context{% if e.In %}, {{ e.In.GoQualifiedModel() }}{% end %}) ({% if e.Out %}{{ e.Out.GoQualifiedModel() }}, {% end %}error)
    {% end %}
}

type {{ repo.StructName() }} struct {
	logger logger.Slog
	client postgres.Client
}

func New{{ repo.StructName() }}(logger logger.Slog, conf postgres.Config) (*{{ repo.StructName() }}, error) {
	client, err := postgres.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &{{ repo.StructName() }}{
		logger: logger,
		client: client,
	}, nil
}