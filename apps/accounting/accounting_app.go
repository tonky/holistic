package accounting

import (
	"context"
	"fmt"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"

	"github.com/samber/do/v2"
)

func (a *App) ReadOrder(ctx context.Context, arg food.OrderID) (accounting.Order, error) {
	return do.MustInvokeAs[AccountOrdersRepoReader](a.deps).ReadOrderByFoodID(ctx, arg)
}

func (a *App) FoodOrderProcessor(ctx context.Context, in food.Order) error {
	a.logger.Info("AccountingApp.FoodOrderProcessor got: ", in)

	repo := do.MustInvokeAs[AccountOrdersRepoReader](a.deps)

	_, err := repo.SaveOrder(ctx, NewOrder{
		Order:   in,
		Content: in.Content,
		Cost:    10,
		IsPaid:  true,
	})

	return err
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
