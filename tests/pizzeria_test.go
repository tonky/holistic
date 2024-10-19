package tests

import (
	"context"
	"testing"
	"tonky/holistic/gen/clients"
	"tonky/holistic/gen/domain/food"
	"tonky/holistic/gen/services/pizzeria"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
)

func init() {
	go func() {
		injector := do.New()
		do.Provide(injector, pizzeria.NewConfig)

		pizzeria := pizzeria.NewPizzeria(injector)
		pizzeria.Start()
	}()
}

func TestPizzeriaCRD(t *testing.T) {
	oid, err := food.NewOrderID("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	injector := do.New()
	do.Provide(injector, pizzeria.NewConfig)
	port := do.MustInvoke[*pizzeria.Config](injector).Port

	conf := clients.Config{
		Host: "localhost",
		Port: port,
	}

	pc := clients.NewPizzeria(conf)

	newOrder := food.Order{
		ID:      oid,
		Content: "new order",
	}

	createdOrder, err := pc.CreateOrder(context.TODO(), newOrder)
	require.NoError(t, err)
	require.Equal(t, newOrder, createdOrder)

	order, err := pc.ReadOrder(context.TODO(), createdOrder.ID)
	require.NoError(t, err)

	require.Equal(t, order.ID, createdOrder.ID)
}
