package pricing

import (
	"context"
	"tonky/holistic/domain/food"
	"tonky/holistic/domain/pricing"
)

func (a App) ReadOrder(ctx context.Context, in food.OrderID) (pricing.OrderPrice, error) {
	return a.Deps.OrdererRepo.ReadOrderPricingByID(ctx, in)
}
