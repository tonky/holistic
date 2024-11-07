package accounting

import (
	"context"
	"fmt"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"

	"github.com/samber/do/v2"
)

func (a *App) ReadOrder(ctx context.Context, arg food.OrderID) (accounting.Order, error) {
	return do.MustInvokeAs[OrdererRepository](a.Deps).ReadOrderByFoodID(ctx, arg)
}

func (a *App) FoodOrderUpdatedProcessor(ctx context.Context, in food.Order) error {
	a.Logger.Info("AccountingApp.FoodOrderUpdatedProcessor got: ", in)

	if !in.IsFinal {
		return nil
	}

	orderPrice, err := a.PricingClient.ReadOrder(ctx, in.ID)
	if err != nil {
		return err
	}

	paidOrder := accounting.Order{
		ID:   in.ID,
		Cost: orderPrice.Cost,
	}

	_, errSave := a.OrdererRepo.SaveFinishedOrder(ctx, paidOrder)
	if errSave != nil {
		return errSave
	}

	return a.AccountingOrderPaidProducer.ProduceAccountingOrderPaid(ctx, paidOrder)
}

func (a *App) FoodOrderProcessorErrHandler(errs chan error) {
	for err := range errs {
		fmt.Println("AccountingApp.foodOrderProcessorErrHandler got error: ", err)
	}
}

type NewOrder struct {
	Order   food.Order
	Content string
	Cost    int
	IsPaid  bool
}
