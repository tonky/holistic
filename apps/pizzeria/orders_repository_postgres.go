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

	err := r.client.PgxConn.QueryRow(context.Background(), "select id, content from orders where id=$1", req.ID).Scan(&id, &content)
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

	_, err := r.client.PgxConn.Exec(ctx, query, args)

	if err != nil {
		return food.Order{}, fmt.Errorf("unable to insert row: %w", err)
	}

	return food.Order{ID: food.OrderID{id}, Content: req.Content}, nil
}

type NewOrder struct {
	Content string
}

/*
func NewPostgresOrdererRepository(deps do.Injector) (*PostgresOrderer, error) {
	return &PostgresOrderer{
		logger: *do.MustInvoke[*logger.SlogLogger](deps),
		client: *do.MustInvoke[*postgres.Client](deps),
	}, nil
}

*/
