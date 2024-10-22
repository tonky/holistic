package {{ service_name }}

import (
	"context"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra"
	"tonky/holistic/infra/logger"
)

type {{ infra.Name}}Repository interface {
    {% for io in infra.InOut %}
    {{ io.Name }}(context.Context{% if io.In %}, {{ io.In.Typ }}{% end %}) ({{ io.Out.Typ}}, error)
    {% end %}
}

type {{ cap(infra.Typ) }}{{ infra.Name }} struct {
	logger logger.SlogLogger
	client infra.{{ cap(infra.Typ) }}Client
}