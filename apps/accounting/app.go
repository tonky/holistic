package accounting

import (
	"context"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"
)

func (a *App) ReadOrder(ctx context.Context, arg food.OrderID) (accounting.Order, error) {
	return a.Deps.OrdererRepo.ReadOrderByFoodID(ctx, arg)
}

func (a *App) FoodOrderUpdatedProcessor(ctx context.Context, in food.Order) error {
	a.Logger.Info("AccountingApp.FoodOrderUpdatedProcessor got: ", in)

	if !in.IsFinal {
		return nil
	}

	orderPrice, err := a.Clients.PricingClient.ReadOrder(ctx, in.ID)
	if err != nil {
		return err
	}

	paidOrder := accounting.Order{
		ID:   in.ID,
		Cost: orderPrice.Cost,
	}

	_, errSave := a.Deps.OrdererRepo.SaveFinishedOrder(ctx, paidOrder)
	if errSave != nil {
		return errSave
	}

	return a.Deps.AccountingOrderPaidProducer.ProduceAccountingOrderPaid(ctx, paidOrder)
}

type NewOrder struct {
	Order   food.Order
	Content string
	Cost    int
	IsPaid  bool
}
