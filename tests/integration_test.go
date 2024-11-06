package tests

import (
	"context"
	"testing"
	"time"
	appAcc "tonky/holistic/apps/accounting"
	app_piz "tonky/holistic/apps/pizzeria"
	"tonky/holistic/clients"
	"tonky/holistic/clients/accountingClient"
	"tonky/holistic/clients/pizzeriaClient"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/kafkaConsumer"
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/postgres"
	svc_acc "tonky/holistic/services/accounting"
	svc_piz "tonky/holistic/services/pizzeria"
	svcPricing "tonky/holistic/services/pricing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
)

func startServices() do.Injector {
	l := logger.Slog{}

	pricing, err := svcPricing.NewFromEnv()
	if err != nil {
		panic(err)
	}

	go pricing.Start()

	time.Sleep(50 * time.Millisecond)

	injector := do.New()
	kc := kafka.Config{Brokers: []string{"localhost:19092"}}
	kfc, err := kafkaConsumer.NewFoodOrderCreatedConsumer(l, kc)
	if err != nil {
		panic(err)
	}

	do.ProvideValue(injector, &kc)
	do.ProvideValue(injector, &svc_piz.Config{Port: 1236})
	do.ProvideValue(injector, &svc_acc.Config{Port: 1235})

	do.ProvideValue(injector, &l)

	pgConf := postgres.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
	}

	po, err := app_piz.NewPostgresOrderer(l, pgConf)
	if err != nil {
		panic(err)
	}

	do.ProvideValue(injector, *po)

	kpoc, err := kafkaProducer.NewFoodOrderCreatedProducer(l, kc)
	if err != nil {
		panic(err)
	}

	kpou, err := kafkaProducer.NewFoodOrderUpdatedProducer(l, kc)
	if err != nil {
		panic(err)
	}

	do.ProvideValue(injector, kpoc)
	do.ProvideValue(injector, kpou)
	do.ProvideValue(injector, kfc)

	pizzeria, err := svc_piz.New(injector)
	if err != nil {
		panic(err)
	}

	go pizzeria.Start()

	kaop, err := kafkaProducer.NewAccountingOrderPaidProducer(l, kc)
	if err != nil {
		panic(err)
	}

	do.ProvideValue(injector, kaop)

	accounting, err := svc_acc.NewFromEnv()
	if err != nil {
		panic(err)
	}

	apgo := do.MustInvoke[*appAcc.PostgresOrderer](accounting.Deps())
	do.ProvideValue(injector, apgo)

	do.ProvideValue(injector, accounting.Config())
	go accounting.Start()

	logger.Slog{}.Info("Test init() - done! Services started")

	return injector
}

func TestOrderThroughKafka(t *testing.T) {
	injector := startServices()

	time.Sleep(100 * time.Millisecond)

	accountingConfig := do.MustInvoke[svc_acc.Config](injector)
	pizzeriaConfig := do.MustInvoke[*svc_piz.Config](injector)

	pizClientConf := clients.Config{Host: "localhost", Port: pizzeriaConfig.Port}
	do.ProvideValue(injector, &pizClientConf)

	pc := pizzeriaClient.New(pizClientConf)

	newOrder := svc_piz.NewOrder{
		Content: "new order",
	}

	createdOrder, err := pc.CreateOrder(context.TODO(), newOrder)
	require.NoError(t, err)
	require.Equal(t, newOrder.Content, createdOrder.Content)

	order, err := pc.ReadOrder(context.TODO(), createdOrder.ID)
	require.NoError(t, err)

	require.Equal(t, order.ID, createdOrder.ID)

	uo := svc_piz.UpdateOrder{
		ID:      createdOrder.ID,
		Content: "updated content",
		IsFinal: true,
	}

	updatedOrder, err := pc.UpdateOrder(context.TODO(), uo)
	require.NoError(t, err)
	require.Equal(t, uo.Content, updatedOrder.Content)
	require.Equal(t, uo.IsFinal, updatedOrder.IsFinal)

	time.Sleep(50 * time.Millisecond)

	ac := accountingClient.New(clients.Config{Host: "localhost", Port: accountingConfig.Port})

	accountingOrder, err := ac.ReadOrder(context.TODO(), createdOrder.ID)
	require.NoError(t, err)

	require.Equal(t, createdOrder.ID, accountingOrder.ID)
	require.Equal(t, accountingOrder.Cost, 5)

	accDB := do.MustInvokeAs[appAcc.OrdererRepository](injector)
	ao, err := accDB.ReadOrderByFoodID(context.TODO(), createdOrder.ID)
	require.NoError(t, err)
	require.Equal(t, ao.Cost, 5)
}
