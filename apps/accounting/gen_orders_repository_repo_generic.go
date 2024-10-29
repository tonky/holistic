package accounting

import (
	"context"

    "tonky/holistic/domain/food"
    "tonky/holistic/domain/accounting"

	"tonky/holistic/infra/logger"
)

// compile-time check to make sure interface is implemented
var _ AccountOrdersRepoReader = new(OrdersRepository)

type AccountOrdersRepoReader interface {
	ReadOrderByFoodID(context.Context, food.OrderID) (accounting.Order, error)
}

type OrdersRepository struct {
	logger logger.Slog
}

func NewOrdersRepository(logger logger.Slog) *OrdersRepository {
	return &OrdersRepository{
        logger: logger,
	}
}
