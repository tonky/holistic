// AUTOGENERATED! DO NOT EDIT.
package {{ service.Name }}Client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
    {% for ci in service.AbsImports(ctx) %}
	"{{ ci }}"
    {% end %}
	"{{ modulePath}}/clients"
	svc "{{ modulePath }}/services/{{ service.Name }}"
)

type I{{ cap(service.Name) }}Client interface {
{% for h in service.Endpoints %}
	{{ h.Name }}(context.Context, {{ h.In.GoStructModel(ctx) }}) ({{ h.Out.GoStructModel(ctx) }}, error)
{% end %}
}

func New(config clients.Config) {{ cap(service.Name) }}Client {
	return {{ cap(service.Name) }}Client {
		config: config,
		logger: {{ service.Logger.Model.Package() }}.Default(),
	}
}

func NewFromEnv() {{ cap(service.Name) }}Client {
	svcConf := svc.MustEnvConfig()

	return {{ cap(service.Name) }}Client {
		config: clients.Config{Host: "http://localhost", Port: svcConf.Port},
        logger: {{ service.Logger.Model.Package() }}.NewFromConfig(svcConf.Logger),
	}
}

type {{ cap(service.Name) }}Client struct {
	config clients.Config
    logger {{ service.Logger.Interface.GoQualifiedModel() }}
}

{% for h in service.Endpoints %}
func (c {{ cap(service.Name) }}Client) {{ h.Name }}(ctx context.Context, arg {{ h.In.GoStructModel(ctx) }}) ({{ h.Out.GoStructModel(ctx) }}, error) {
	var reply {{ h.Out.GoStructModel(ctx) }}

	jsonBody, err := json.Marshal(arg)
	if err != nil { return reply, err}

	bodyReader := bytes.NewReader(jsonBody)

 	requestURL := fmt.Sprintf("%s/%s", c.config.ServerAddress(), "{{ h.Name }}")

	c.logger.Debug("{{ cap(service.Name) }}Client.{{ h.Name }}()", "requestURL", requestURL, "newRecipe", arg)

 	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil { return reply, err }

	res, err := http.DefaultClient.Do(req)
	if err != nil { return reply, err }

	resBody, err := io.ReadAll(res.Body)
	if err != nil { return reply, err }

	if err := json.Unmarshal(resBody, &reply); err != nil { return reply, err }

	return reply, nil
}

{% end %}