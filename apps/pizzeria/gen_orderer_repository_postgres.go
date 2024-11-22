// AUTOGENERATED! DO NOT EDIT.
package pizzeria

import (
	"context"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/postgres"
	"tonky/holistic/domain/food"
)

var _ OrdererRepository = new(PostgresOrderer)

type OrdererRepository interface {
    ReadOrderByID(context.Context, food.OrderID) (food.Order, error)
    SaveOrder(context.Context, NewOrder) (food.Order, error)
    UpdateOrder(context.Context, UpdateOrder) (food.Order, error)
}

type PostgresOrderer struct {
	logger logger.ILogger
	client postgres.Client
}

func NewPostgresOrderer(l logger.ILogger, conf postgres.Config) (*PostgresOrderer, error) {
	client, err := postgres.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &PostgresOrderer{
		logger: l,
		client: client,
	}, nil
}
