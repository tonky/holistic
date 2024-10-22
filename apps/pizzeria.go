package apps

import (
	"context"
	"tonky/holistic/domain/food"

	"github.com/samber/do/v2"
)

type PizzeriaApp struct {
	deps do.Injector
}

func NewPizzeria(deps do.Injector) PizzeriaApp {
	return PizzeriaApp{deps: deps}
}

func (a PizzeriaApp) ReadOrder(ctx context.Context, id food.OrderID) (food.Order, error) {
	res := food.Order{ID: id}
	return res, nil
}

func (a PizzeriaApp) CreateOrder(ctx context.Context, in food.Order) (food.Order, error) {
	// var res food.Order

	return in, nil
}
