package accounting

import (
	"context"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"
)

func (a OrdersRepository) ReadOrderByFoodID(ctx context.Context, in food.OrderID) (accounting.Order, error) {
	var out accounting.Order

	a.logger.Info("AccountingOrdersRepository.ReadOrderByID", in)

	return out, nil
}
