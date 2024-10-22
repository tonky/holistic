package apps

import (
	"context"
	"log/slog"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra"
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

type PizzeriaApp struct {
	deps     do.Injector
	postgres *infra.PostgresClient
	logger   *logger.SlogLogger
}

func NewPizzeria(deps do.Injector) PizzeriaApp {
	pgConf := do.MustInvoke[infra.PostgresConfig](deps)
	pgc, err := infra.NewPosgresClient(pgConf)
	if err != nil {
		panic(err)
	}

	return PizzeriaApp{
		deps:     deps,
		logger:   do.MustInvoke[*logger.SlogLogger](deps),
		postgres: pgc,
	}
}

func (app PizzeriaApp) ReadOrder(ctx context.Context, id food.OrderID) (food.Order, error) {
	app.logger.Info("ReadOrder pgxConn", slog.Any("pg conn", app.postgres))

	res := food.Order{ID: id}
	return res, nil
}

func (app PizzeriaApp) CreateOrder(ctx context.Context, in food.Order) (food.Order, error) {
	// var res food.Order

	return in, nil
}
