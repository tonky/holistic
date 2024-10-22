// AUTOGENERATED! DO NOT EDIT.
package pizzeria

import (
    "context"
    "fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

    // "github.com/go-playground/validator/v10"
    "github.com/samber/do/v2"

	"tonky/holistic/domain/food"
	"tonky/holistic/apps/pizzeria"
)

type Pizzeria struct {
    config Config
    deps do.Injector
}

func (h Pizzeria) ReadOrder(arg food.OrderID, reply *food.Order) error {
    application, appErr := pizzeria.NewApp(h.deps)
    if appErr != nil {
        return appErr
    }

    res, err := application.ReadOrder(context.TODO(), arg)
    if err != nil {
        return err
    }

    *reply = res

    return nil
}

func (h Pizzeria) CreateOrder(arg food.Order, reply *food.Order) error {
    application, appErr := pizzeria.NewApp(h.deps)
    if appErr != nil {
        return appErr
    }

    res, err := application.CreateOrder(context.TODO(), arg)
    if err != nil {
        return err
    }

    *reply = res

    return nil
}


func NewPizzeria(dependencies do.Injector) ServiceStarter {
	cfg := do.MustInvoke[*Config](dependencies)

    handlers := Pizzeria{deps: dependencies, config: *cfg}

    return handlers
}

func (h Pizzeria) Start() error {
	port := h.config.Port

    fmt.Printf(">> pizzeria.Start() config: %+v\n", h.config)

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
