package clients


import (
	"context"
	"log"
	"net/rpc"
	"tonky/holistic/gen/domain/food"
    // generate imports
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
func (c {{ cap(service_name) }}Client) {{ h.FuncName() }}(ctx context.Context, req {{ h.In }}) ({{ h.Out.ok }}, error) {
	client, err := rpc.DialHTTP("tcp", c.config.ServerAddress())
	if err != nil {
		log.Fatal("dialing error:", err)
	}

	var reply food.Order

	err = client.Call("Pizzeria.{{ h.FuncName() }}", req, &reply)
	if err != nil {
		log.Fatal("server call error:", err)
	}

	return reply, nil
}

{% end %}