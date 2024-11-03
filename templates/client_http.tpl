// AUTOGENERATED! DO NOT EDIT.
package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	{% for imp in client_relative_imports %}
	{{ imp.Alias}} "tonky/holistic/{{ imp.RelPath }}"
	{% end %}
)

func New{{ cap(service_name) }}(config Config) {{ cap(service_name) }}Client {
	return {{ cap(service_name) }}Client{
		config: config,
	}
}

type {{ cap(service_name) }}Client struct {
	config Config
}

{% for h in handlers %}
func (c {{ cap(service_name) }}Client) {{ h.FuncName() }}(ctx context.Context, arg {{ h.In.ServiceModel() }}) ({{ h.Out.ok }}, error) {
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