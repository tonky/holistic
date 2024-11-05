package accounting

import (
	"context"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"
)

func (a OrdersRepository) ReadOrderByFoodID(ctx context.Context, in food.OrderID) (accounting.Order, error) {
	a.logger.Info("AccountingOrdersRepository.ReadOrderByID", in)

	return accounting.Order{ID: in, Cost: 5}, nil
}

func (a OrdersRepository) SaveOrder(ctx context.Context, in NewOrder) (accounting.Order, error) {
	var out accounting.Order

	a.logger.Info("AccountingOrdersRepository.SaveAccountingOrder", in)

	return out, nil
}
