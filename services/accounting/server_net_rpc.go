// AUTOGENERATED! DO NOT EDIT.
package accounting

import (
    "context"
    "fmt"
	"log"
	"net"
	"net/rpc"

    // "github.com/go-playground/validator/v10"
    "github.com/samber/do/v2"

	"tonky/holistic/infra/kafkaConsumer"

	"tonky/holistic/domain/food"
	"tonky/holistic/domain/accounting"
	app "tonky/holistic/apps/accounting"
	"tonky/holistic/infra/logger"
)

type Accounting struct {
    config Config
    deps do.Injector
    app app.App
}

func (h Accounting) ReadOrder(arg food.OrderID, reply *accounting.Order) error {
    res, err := h.app.ReadOrder(context.TODO(), arg)
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

    handlers := Accounting{deps: dependencies, config: *cfg, app: *application}

    return handlers, nil
}

func (h Accounting) Start() error {
	port := h.config.Port

    fmt.Printf(">> accounting.Start() config: %+v\n", h.config)

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

	FoodOrderUpdatedConsumer, err := kafkaConsumer.NewFoodOrderUpdatedConsumer(l, cfg.App.Kafka)
    if err != nil {
        return nil, err
    }

	do.ProvideValue(deps, FoodOrderUpdatedConsumer)
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

    handlers := Accounting{deps: deps, config: cfg, app: *application}

    return handlers, nil
}

func (h Accounting) Config() Config {
    return h.config
}

func (h Accounting) Deps() do.Injector {
    return h.deps
}


type ServiceStarter interface {
    Start() error
    Config() Config
    Deps() do.Injector
}
