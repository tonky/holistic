package accounting

import (
	"context"
	"fmt"
	"time"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"

	"github.com/jackc/pgx/v5"
)

func (a PostgresOrderer) ReadOrderByFoodID(ctx context.Context, in food.OrderID) (accounting.Order, error) {
	a.logger.Info("accounting.PostgresOrderer.ReadOrderByID", in)

	var cost int
	var paid_at time.Time

	err := a.client.PgxConn.QueryRow(context.Background(), "select price, paid_at from accounting_orders where order_id=$1", in).Scan(&cost, &paid_at)
	if err != nil {
		return accounting.Order{}, fmt.Errorf("order query for id %s failed: %v", in.ID.String(), err)
	}

	a.logger.Debug("accounting.PostgresOrderer.ReadOrderByID", "cost", cost, "paid_at", paid_at)

	return accounting.Order{ID: in, Cost: cost}, nil
}

func (a PostgresOrderer) SaveFinishedOrder(ctx context.Context, in accounting.Order) (accounting.Order, error) {
	a.logger.Info("accounting.PostgresOrderer..SaveAccountingOrder", in)

	query := `INSERT INTO accounting_orders (order_id, price, paid_at) VALUES (@id, @price, @time)`
	args := pgx.NamedArgs{"id": in.ID.String(), "price": in.Cost, "time": time.Now()}

	_, err := a.client.PgxConn.Exec(ctx, query, args)

	if err != nil {
		return accounting.Order{}, fmt.Errorf("unable to insert row: %w", err)
	}

	return in, nil
}
