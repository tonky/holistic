package app

import (
	"context"
	"fmt"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra"
	"tonky/holistic/services/pizzeria"

	"github.com/samber/do/v2"
)

type App struct {
	deps do.Injector
}

type Orderer interface {
	ReadOrder(context.Context, food.OrderID) (food.Order, error)
	CreateOrder(context.Context, food.Order) (food.Order, error)
}

func New(deps do.Injector) Orderer {
	cfg := do.MustInvoke[*pizzeria.Config](deps)

	if cfg.ShouldMockApp {
		return NewMock()
	}

	return App{deps: deps}
}

func (a App) ReadOrder(ctx context.Context, req food.OrderID) (food.Order, error) {
	fmt.Println("App.ReadOrder", req)

	return food.Order{ID: req, Content: "idk"}, nil
}

func (a App) CreateOrder(ctx context.Context, req food.Order) (food.Order, error) {
	fmt.Println("App.CrateOrder", req)

	fmt.Println("App.CrateOrder DB", do.MustInvoke[*infra.PostgresClient](a.deps))

	return req, nil
}
