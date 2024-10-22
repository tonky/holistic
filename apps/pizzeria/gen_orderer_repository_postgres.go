package pizzeria

import (
	"context"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra"
	"tonky/holistic/infra/logger"
)

type OrdererRepository interface {
    ReadOrderByID(context.Context, food.OrderID) (food.Order, error)
    SaveOrder(context.Context, NewOrder) (food.Order, error)
}

type PostgresOrderer struct {
	logger logger.SlogLogger
	client infra.PostgresClient
}