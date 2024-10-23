package pizzeria

import (
	"context"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra"
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

func (r PostgresOrderer) ReadOrderByID(ctx context.Context, req food.OrderID) (food.Order, error) {
	return food.Order{}, nil
}

func (r PostgresOrderer) SaveOrder(ctx context.Context, req NewOrder) (food.Order, error) {
	r.logger.Info("PostgresOrderer | CreateOrder", "pg conn", r.client, req)
	return food.Order{}, nil
}

type NewOrder struct {
	Content string
}

func NewPostgresOrderRepository(deps do.Injector) (*PostgresOrderer, error) {
	postgresConf := do.MustInvoke[*infra.PostgresConfig](deps)

	postgresClient, err := infra.NewPostgresClient(*postgresConf)
	if err != nil {
		return nil, err
	}

	return &PostgresOrderer{
		logger: *do.MustInvoke[*logger.SlogLogger](deps),
		client: postgresClient,
	}, nil
}
