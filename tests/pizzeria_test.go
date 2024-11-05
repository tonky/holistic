package tests

import (
	"context"
	"fmt"
	"testing"
	"time"
	app "tonky/holistic/apps/pizzeria"
	"tonky/holistic/clients"
	"tonky/holistic/clients/pizzeriaClient"
	"tonky/holistic/infra/logger"
	svc "tonky/holistic/services/pizzeria"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
)

func initPizzeriaServer() do.Injector {
	injector := do.New()

	do.ProvideValue(injector, &svc.Config{Port: 1234})
	do.ProvideValue(injector, &logger.Slog{})

	do.Provide(injector, app.NewMemoryOrdererRepository)
	do.Provide(injector, app.NewMemoryOrderProducerRepository)

	pizzeria, err := svc.New(injector)
	if err != nil {
		panic(err)
	}

	go pizzeria.Start()

	return injector
}

func TestPizzeriaCRD(t *testing.T) {
	deps := initPizzeriaServer()

	time.Sleep(100 * time.Millisecond)

	conf := clients.Config{Host: "localhost", Port: do.MustInvoke[*svc.Config](deps).Port}

	do.ProvideValue(deps, &conf)

	pc := pizzeriaClient.New(conf)

	newOrder := svc.NewOrder{
		Content: "new order",
	}

	fmt.Println("calling server")
	createdOrder, err := pc.CreateOrder(context.TODO(), newOrder)
	require.NoError(t, err)
	require.Equal(t, newOrder.Content, createdOrder.Content)

	order, err := pc.ReadOrder(context.TODO(), createdOrder.ID)
	require.NoError(t, err)

	require.Equal(t, order.ID, createdOrder.ID)
	// require.False(t, true)
}
