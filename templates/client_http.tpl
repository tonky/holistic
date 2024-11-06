// AUTOGENERATED! DO NOT EDIT.
package {{ service.Name }}Client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"tonky/holistic/clients"
	{% for imp in client_relative_imports %}
	{{ imp.Alias}} "tonky/holistic/{{ imp.RelPath }}"
	{% end %}
)

type I{{ cap(service.Name) }}Client interface {
{% for h in handlers %}
	{{ h.FuncName() }}(context.Context, {{ h.In.ServiceModel() }}) ({{ h.Out.ok }}, error)
{% end %}
}

func New(config clients.Config) {{ cap(service.Name) }}Client {
	return {{ cap(service.Name) }}Client{
		config: config,
	}
}

func NewFromEnv(env string) {{ cap(service.Name) }}Client {
	envConf := clients.ConfigForEnv("{{ service.Name }}", env)

	return {{ cap(service.Name) }}Client{
		config: envConf,
	}
}

type {{ cap(service.Name) }}Client struct {
	config clients.Config
}

{% for h in handlers %}
func (c {{ cap(service.Name) }}Client) {{ h.FuncName() }}(ctx context.Context, arg {{ h.In.ServiceModel() }}) ({{ h.Out.ok }}, error) {
	var reply {{ h.Out.ok }}

	jsonBody, err := json.Marshal(arg)
	if err != nil { return reply, err}

	bodyReader := bytes.NewReader(jsonBody)

 	requestURL := fmt.Sprintf("%s/%s", c.config.ServerAddress(), "{{ h.Method.HttpName() }}{{ cap(h.Name) }}")

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