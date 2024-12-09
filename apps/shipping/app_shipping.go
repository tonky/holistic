package shipping

import (
	"context"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"
	"tonky/holistic/domain/shipping"
)

func (a App) ReadOrder(ctx context.Context, in food.OrderID) (shipping.Order, error) {
	a.Deps.Logger.Info("shipping.App.ReadOrder", in)

	return a.Deps.OrdererRepo.ReadOrderShippingByID(ctx, in)
}

func (a App) AccountingOrderPaidProcessor(ctx context.Context, in accounting.Order) error {
	a.Deps.Logger.Info("shipping.App.AccountingOrderPaidProcessor", in)

	order, err := a.Deps.OrdererRepo.SaveShipping(ctx, shipping.Order{ID: in.ID})
	if err != nil {
		return err
	}

	return a.Deps.ShippingOrderShippedProducer.ProduceShippingOrderShipped(ctx, order)
}
