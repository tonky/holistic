package {{ service.Name }}

import (
	"context"

    {% for i in repo.Imports(service) %}
    {{ i.String(mod) }}
    {% end %}

    {% for dep := range repo.Deps %}
	"tonky/holistic/{{ repo.Deps[dep].FQImport(service) }}"
    {% end %}
)

// compile-time check to make sure interface is implemented
var _ {{ repo.InterfaceName() }} = new({{ repo.Struct }})

type {{ repo.InterfaceName() }} interface {
    {% for m in repo.Methods %}
	{{ m.Name }}(context.Context, {{ m.Arg.Typ }}) ({{ m.Ret.Typ }}, error)
    {% end %}
}

type {{ repo.Struct }} struct {
    {% for dep := range repo.Deps %}
	{{ dep }} {{ repo.Deps[dep].FQModel() }}
    {% end %}
}

func New{{ repo.Struct }}({{ repo.Deps.StructArgsStr() }}) *{{ repo.Struct }} {
	return &{{ repo.Struct }}{
        {% for dep := range repo.Deps %}
        {{ dep }}: {{ dep }},
        {% end %}
	}
}
