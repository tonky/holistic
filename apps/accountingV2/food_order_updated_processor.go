package accountingV2

import (
	"context"
	"log/slog"
	"tonky/holistic/domain/food"
)

func (a App) FoodOrderUpdatedProcessor(ctx context.Context, oder food.Order) error {
	slog.Info("accountingV2.App.foodOrderUpdatedProcessor()")

	return nil
}
