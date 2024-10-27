package accounting

import (
	"context"
	"fmt"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"
)

func (a *App) ReadOrder(ctx context.Context, arg food.OrderID) (accounting.Order, error) {
	res := accounting.Order{}

	return res, nil
}

func (a *App) FoodOrderProcessor(ctx context.Context, in food.Order) error {
	a.logger.Info("AccountingApp.FoodOrderProcessor got: ", in)

	return nil
}

func (a *App) FoodOrderProcessorErrHandler(errs chan error) {
	for err := range errs {
		fmt.Println("AccountingApp.foodOrderProcessorErrHandler got error: ", err)
	}
}
