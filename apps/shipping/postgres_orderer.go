package shipping

import (
	"context"
	"fmt"
	"time"
	"tonky/holistic/domain/food"
	"tonky/holistic/domain/shipping"

	"github.com/jackc/pgx/v5"
)

func (a PostgresOrderer) ReadOrderShippingByID(ctx context.Context, in food.OrderID) (shipping.Order, error) {
	a.logger.Info("shipping.PostgresOrderer.ReadOrderByFoodID", in)

	var shippedAt time.Time

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := a.client.Pool.QueryRow(ctx, "select shipped_at from shipping where order_id=$1", in).Scan(&shippedAt)
	if err != nil {
		return shipping.Order{}, fmt.Errorf("order query for id %s failed: %v", in.ID.String(), err)
	}

	return shipping.Order{ID: in, ShippedAt: shippedAt}, nil
}

func (a PostgresOrderer) SaveShipping(ctx context.Context, in shipping.Order) (shipping.Order, error) {
	a.logger.Info("shipping.PostgresOrderer.SaveShipping", in)

	query := `INSERT INTO shipping (order_id, shipped_at) VALUES (@id, @time)`
	args := pgx.NamedArgs{"id": in.ID.String(), "time": time.Now()}

	_, err := a.client.Pool.Exec(ctx, query, args)

	if err != nil {
		return shipping.Order{}, fmt.Errorf("unable to insert shipping order row: %w", err)
	}

	return in, nil
}
