package accountingV2

import (
	"context"
	"tonky/holistic/domain/food"
)

func (a App) FoodOrderUpdatedProcessor(ctx context.Context, oder food.Order) error {
	a.Logger.Info(">> accountingV2.App.foodOrderUpdatedProcessor()")

	return nil
}
