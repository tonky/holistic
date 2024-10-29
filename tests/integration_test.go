package tests

import (
	"context"
	"testing"
	"time"
	app_acc "tonky/holistic/apps/accounting"
	app_piz "tonky/holistic/apps/pizzeria"
	"tonky/holistic/clients"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/logger"
	svc_acc "tonky/holistic/services/accounting"
	svc_piz "tonky/holistic/services/pizzeria"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
)

func startServices() do.Injector {
	injector := do.New()
	kc := kafka.Config{Brokers: []string{"localhost:19092"}}
	kfc, err := app_acc.NewKafkaFoodOrderConsumer(logger.Slog{}, kc)
	if err != nil {
		panic(err)
	}

	do.ProvideValue(injector, &kc)
	do.ProvideValue(injector, &svc_piz.Config{Port: 1236, Kafka: kc})
	do.ProvideValue(injector, &svc_acc.Config{Port: 1235, Kafka: kc})

	do.ProvideValue(injector, &logger.Slog{})

	do.Provide(injector, app_piz.NewMemoryOrdererRepository)
	do.Provide(injector, app_piz.NewDOKafkaFoodOrderProducer)

	do.ProvideValue(injector, kfc)

	pizzeria, err := svc_piz.NewPizzeria(injector)
	if err != nil {
		panic(err)
	}

	go pizzeria.Start()

	accounting, err := svc_acc.NewAccounting(injector)
	if err != nil {
		panic(err)
	}

	go accounting.Start()

	logger.Slog{}.Info("Test init() - done! Services started")

	return injector
}

func TestOrderThroughKafka(t *testing.T) {
	injector := startServices()

	time.Sleep(500 * time.Millisecond)

	accountingConfig := do.MustInvoke[*svc_acc.Config](injector)
	pizzeriaConfig := do.MustInvoke[*svc_piz.Config](injector)

	conf := clients.Config{Host: "localhost", Port: pizzeriaConfig.Port}

	do.ProvideValue(injector, &conf)

	pc := clients.NewPizzeria(conf)

	newOrder := svc_piz.NewOrder{
		Content: "new order",
	}

	createdOrder, err := pc.CreateOrder(context.TODO(), newOrder)
	require.NoError(t, err)
	require.Equal(t, newOrder.Content, createdOrder.Content)

	order, err := pc.ReadOrder(context.TODO(), createdOrder.ID)
	require.NoError(t, err)

	require.Equal(t, order.ID, createdOrder.ID)

	time.Sleep(1000 * time.Millisecond)

	ac := clients.NewAccounting(clients.Config{Host: "localhost", Port: accountingConfig.Port})

	accountingOrder, err := ac.ReadOrder(context.TODO(), createdOrder.ID)
	require.NoError(t, err)

	require.Equal(t, accountingOrder.ID, createdOrder.ID)

	require.Equal(t, accountingOrder.IsPaid, true)
	require.Equal(t, accountingOrder.Cost, 10)

}