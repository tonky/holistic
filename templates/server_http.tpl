package service

import (
    "context"
    "fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

    // "github.com/go-playground/validator/v10"
    "github.com/samber/do/v2"

	"tonky/holistic/gen/domain/food"
	"tonky/holistic/gen/services/{{ service_name }}"
	"tonky/holistic/gen/services/{{ service_name }}/app"
)

type {{ cap(service_name) }} struct {
    deps do.Injector
}

{% for h in handlers %}
func (h {{ cap(service_name) }}) {{h.FuncName()}}(arg {{ h.In }}, reply *{{ h.Out.ok }}) error {
    application := app.New(h.deps)

    res, err := application.{{h.FuncName()}}(context.TODO(), arg)
    if err != nil {
        return err
    }

    *reply = res

    return nil
}

{% end %}

func New{{ cap(service_name) }}(dependencies do.Injector) ServiceStarter {
    handlers := {{ cap(service_name) }}{deps: dependencies}
    return handlers
}

func (h {{ cap(service_name) }}) Start() error {
	cfg := do.MustInvoke[*{{ service_name}}.Config](h.deps)
	port := cfg.Port

    fmt.Printf(">> {{ service_name }}.Start() config: %+v\n", cfg)

	rpc.Register(h)
	rpc.HandleHTTP()

	fmt.Println(">> starging server on port ", port)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatal("listen error:", err)
    }

    return http.Serve(l, nil)
}

type ServiceStarter interface {
    Start() error
}
