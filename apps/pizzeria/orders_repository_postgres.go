package pizzeria

import (
	"context"
	"fmt"
	"tonky/holistic/domain/food"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r PostgresOrderer) ReadOrderByID(ctx context.Context, req food.OrderID) (food.Order, error) {
	r.logger.Info("PostgresOrderer | ReadOrderByID", req)

	var id, content string

	err := r.client.Pool.QueryRow(context.Background(), "select id, content from orders where id=$1", req.ID).Scan(&id, &content)
	if err != nil {
		return food.Order{}, fmt.Errorf("order query for id %s failed: %v", req.ID.String(), err)
	}

	uid, err := food.NewOrderID(id)
	if err != nil {
		return food.Order{}, err
	}

	return food.Order{ID: uid, Content: content}, nil
}

func (r PostgresOrderer) SaveOrder(ctx context.Context, req NewOrder) (food.Order, error) {
	r.logger.Info("PostgresOrderer | CreateOrder", "pg conn", r.client, req)
	id := uuid.New()

	query := `INSERT INTO orders (id, content) VALUES (@id, @content)`
	args := pgx.NamedArgs{"id": id.String(), "content": req.Content}

	_, err := r.client.Pool.Exec(ctx, query, args)

	if err != nil {
		return food.Order{}, fmt.Errorf("unable to insert row: %w", err)
	}

	return food.Order{ID: food.OrderID{ID: id}, Content: req.Content}, nil
}

func (r PostgresOrderer) UpdateOrder(ctx context.Context, req UpdateOrder) (food.Order, error) {
	r.logger.Info("PostgresOrderer | UpdateOrder", req)

	query := `UPDATE orders SET content = @content, is_final = @is_final WHERE id = @id`
	args := pgx.NamedArgs{"id": req.OrderID.String(), "content": req.Content, "is_final": req.IsFinal}

	_, err := r.client.Pool.Exec(ctx, query, args)

	if err != nil {
		return food.Order{}, fmt.Errorf("unable to insert row: %w", err)
	}

	return food.Order{ID: req.OrderID, Content: req.Content, IsFinal: req.IsFinal}, nil
}
