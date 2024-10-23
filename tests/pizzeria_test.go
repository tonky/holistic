package tests

import (
	"context"
	"testing"
	"time"
	app "tonky/holistic/apps/pizzeria"
	"tonky/holistic/clients"
	svc "tonky/holistic/services/pizzeria"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
)

func init() {
	go func() {
		injector := do.New()
		do.Provide(injector, svc.NewConfig)
		do.Provide(injector, app.NewMemoryOrdererRepository)

		pizzeria := svc.NewPizzeria(injector)
		pizzeria.Start()
	}()

	time.Sleep(100 * time.Millisecond)
}

func TestPizzeriaCRD(t *testing.T) {
	injector := do.New()
	do.Provide(injector, svc.NewConfig)
	port := do.MustInvoke[*svc.Config](injector).Port

	conf := clients.Config{
		Host: "localhost",
		Port: port,
	}

	pc := clients.NewPizzeria(conf)

	newOrder := svc.NewOrder{
		Content: "new order",
	}

	createdOrder, err := pc.CreateOrder(context.TODO(), newOrder)
	require.NoError(t, err)
	require.Equal(t, newOrder.Content, createdOrder.Content)

	order, err := pc.ReadOrder(context.TODO(), createdOrder.ID)
	require.NoError(t, err)

	require.Equal(t, order.ID, createdOrder.ID)
	// require.False(t, true)
}
