package pricing

import (
	"context"
	"tonky/holistic/domain/food"
	"tonky/holistic/domain/pricing"
)

func (p PostgresOrderer) ReadOrderPricingByID(ctx context.Context, in food.OrderID) (pricing.OrderPrice, error) {
	p.logger.Info("pricing.PostgresOrderer.ReadOrderPricingByID", in)

	return pricing.OrderPrice{ID: in, Cost: 5}, nil
}

func (p PostgresOrderer) SaveOrderPrice(ctx context.Context, in pricing.OrderPrice) (pricing.OrderPrice, error) {
	p.logger.Info("pricing.PostgresOrderer.SaveOrderPrice", in)

	return in, nil
}
