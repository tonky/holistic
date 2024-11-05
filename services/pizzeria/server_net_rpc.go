// AUTOGENERATED! DO NOT EDIT.
package pizzeria

import (
    "context"
    "fmt"
	"log"
	"net"
	"net/rpc"

    // "github.com/go-playground/validator/v10"
    "github.com/samber/do/v2"

	"tonky/holistic/infra/kafkaProducer"

	"tonky/holistic/domain/food"
	app "tonky/holistic/apps/pizzeria"
	"tonky/holistic/infra/logger"
)

type Pizzeria struct {
    config Config
    deps do.Injector
    app app.App
}

func (h Pizzeria) ReadOrder(arg food.OrderID, reply *food.Order) error {
    res, err := h.app.ReadOrder(context.TODO(), arg)
    if err != nil {
        return err
    }

    *reply = res

    return nil
}

func (h Pizzeria) CreateOrder(arg NewOrder, reply *food.Order) error {
    res, err := h.app.CreateOrder(context.TODO(), arg.ToApp())
    if err != nil {
        return err
    }

    *reply = res

    return nil
}

func (h Pizzeria) UpdateOrder(arg UpdateOrder, reply *food.Order) error {
    res, err := h.app.UpdateOrder(context.TODO(), arg.ToApp())
    if err != nil {
        return err
    }

    *reply = res

    return nil
}


func New(dependencies do.Injector) (ServiceStarter, error) {
	cfg := do.MustInvoke[*Config](dependencies)

    application, appErr := app.NewApp(dependencies)
    if appErr != nil {
        return nil, appErr
    }

    handlers := Pizzeria{deps: dependencies, config: *cfg, app: *application}

    return handlers, nil
}

func (h Pizzeria) Start() error {
	port := h.config.Port

    fmt.Printf(">> pizzeria.Start() config: %+v\n", h.config)

	server := rpc.NewServer()
	server.Register(h)

	fmt.Println(">> starging server on port ", port)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatal("listen error:", err)
    }

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			server.ServeConn(conn)
		}()
	}
}

func NewFromEnv() (ServiceStarter, error) {
	cfg, err := NewEnvConfig()
    if err != nil {
        return nil, err
    }

    deps := do.New()

	do.ProvideValue(deps, cfg)

    l := logger.Slog{}
	do.ProvideValue(deps, &l)

	ordererRepo, err := app.NewPostgresOrderer(l, cfg.App.PostgresOrderer)
    if err != nil {
        return nil, err
    }

	do.ProvideValue(deps, ordererRepo)
	FoodOrderCreatedProducer, err := kafkaProducer.NewFoodOrderCreatedProducer(l, cfg.App.Kafka)
    if err != nil {
        return nil, err
    }

	do.ProvideValue(deps, FoodOrderCreatedProducer)
/*
	ocp, err := kafkaProducer.NewFoodOrderCreatedProducer(l, cfg.App.Kafka)
    if err != nil {
        return nil, err
    }

	oup, err := kafkaProducer.NewFoodOrderUpdatedProducer(l, cfg.App.Kafka)
    if err != nil {
        return nil, err
    }

	do.ProvideValue(deps, ocp)
	do.ProvideValue(deps, oup)
*/
    application, appErr := app.NewApp(deps)
    if appErr != nil {
        return nil, appErr
    }

    handlers := Pizzeria{deps: deps, config: cfg, app: *application}

    return handlers, nil
}

func (h Pizzeria) Config() Config {
    return h.config
}

func (h Pizzeria) Deps() do.Injector {
    return h.deps
}


type ServiceStarter interface {
    Start() error
    Config() Config
    Deps() do.Injector
}
