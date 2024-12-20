// AUTOGENERATED! DO NOT EDIT.
package {{ service.Name }}

import (
	"context"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/postgres"
    {% for ci in repo.Imports(service) %}
	{{ ci.Alias }}"tonky/holistic/{{ ci.RelPath }}"
    {% end %}
)

var _ {{ repo.InterfaceName() }} = new({{ repo.StructName() }})

type {{ repo.InterfaceName() }} interface {
    {% for io in repo.Methods %}
    {{ io.Name }}(context.Context{% if io.Arg %}, {{ io.Arg.Typ }}{% end %}) ({% if io.Ret %}{{ io.Ret.Typ }}, {% end %}error)
    {% end %}
}

type {{ repo.StructName() }} struct {
	logger logger.ILogger
	client postgres.Client
}

func New{{ repo.StructName() }}(l logger.ILogger, conf postgres.Config) (*{{ repo.StructName() }}, error) {
	client, err := postgres.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &{{ repo.StructName() }}{
		logger: l,
		client: client,
	}, nil
}
