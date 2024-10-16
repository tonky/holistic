package main

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

    reply = &res

    return nil
}

{% end %}

func main() {
    dependencies := do.New()
    // do.ProvideValue(dependencies, &Config{ Port: 4242, })

    // car, err := do.Invoke[*Car](i)
    // if err != nil { log.Fatal(err.Error()) }
    // car.Start()

    handlers := {{ cap(service_name) }}{deps: dependencies}
    
    rpc.Register(handlers)

    rpc.HandleHTTP()

    fmt.Println(">> starging server on port 1234")

    l, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("listen error:", err)
    }

    http.Serve(l, nil)
}

