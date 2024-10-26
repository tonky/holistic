package tests

import (
	"context"
	"testing"
	"time"
	app "tonky/holistic/apps/pizzeria"
	"tonky/holistic/clients"
	"tonky/holistic/infra/logger"
	svc "tonky/holistic/services/pizzeria"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
)

func init() {
	go func() {
		injector := do.New()

		do.ProvideValue(injector, &svc.Config{Port: 1234})
		do.ProvideValue(injector, &logger.Slog{})

		do.Provide(injector, app.NewMemoryOrdererRepository)
		do.Provide(injector, app.NewMemoryOrderProducerRepository)

		pizzeria, err := svc.NewPizzeria(injector)
		if err != nil {
			panic(err)
		}

		pizzeria.Start()
	}()

	time.Sleep(100 * time.Millisecond)
}

func TestPizzeriaCRD(t *testing.T) {
	injector := do.New()

	conf := clients.Config{Host: "localhost", Port: 1234}

	do.ProvideValue(injector, &conf)

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
