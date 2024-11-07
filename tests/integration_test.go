package tests

import (
	"context"
	"os"
	"testing"
	"time"
	appAcc "tonky/holistic/apps/accounting"
	app_piz "tonky/holistic/apps/pizzeria"
	appShipping "tonky/holistic/apps/shipping"
	"tonky/holistic/clients"
	"tonky/holistic/clients/accountingClient"
	"tonky/holistic/clients/pizzeriaClient"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/kafkaConsumer"
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/infra/logger"
	svc_acc "tonky/holistic/services/accounting"
	svc_piz "tonky/holistic/services/pizzeria"
	svcPricing "tonky/holistic/services/pricing"
	svcShipping "tonky/holistic/services/shipping"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
)

func startServices() do.Injector {
	os.Setenv("PIZZERIA_PORT", "1234")
	os.Setenv("PIZZERIA_APP_POSTGRESORDERER_PORT", "5432")
	os.Setenv("PIZZERIA_APP_KAFKA_GROUP_ID", "test-order-through-kafka-pizzeria")

	os.Setenv("PRICING_PORT", "1235")
	os.Setenv("PRICING_APP_POSTGRESORDERER_PORT", "5433")

	os.Setenv("ACCOUNTING_PORT", "1236")
	os.Setenv("ACCOUNTING_APP_POSTGRESORDERER_PORT", "5434")
	os.Setenv("ACCOUNTING_APP_KAFKA_GROUP_ID", "test-order-through-kafka-accounting")

	os.Setenv("SHIPPING_PORT", "1237")
	os.Setenv("SHIPPING_APP_POSTGRESORDERER_PORT", "5435")
	os.Setenv("SHIPPING_APP_KAFKA_GROUP_ID", "test-order-through-kafka-shipping")

	injector := do.New()

	l := logger.Slog{}

	shipping, err := svcShipping.NewFromEnv()
	if err != nil {
		panic(err)
	}

	go shipping.Start()

	// sdb := do.MustInvokeAs[*appShipping.PostgresOrderer](shipping.Deps())
	// do.ProvideValue(injector, sdb)
	do.ProvideValue(injector, shipping.Deps().OrdererRepo)

	pricing, err := svcPricing.NewFromEnv()
	if err != nil {
		panic(err)
	}

	go pricing.Start()

	time.Sleep(50 * time.Millisecond)

	kc := kafka.Config{Brokers: []string{"localhost:19092"}}
	kfc, err := kafkaConsumer.NewFoodOrderCreatedConsumer(l, kc)
	if err != nil {
		panic(err)
	}

	do.ProvideValue(injector, &kc)
	do.ProvideValue(injector, &l)

	po, err := app_piz.NewPostgresOrderer(l, svc_piz.MustEnvConfig().App.PostgresOrderer)
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

	accounting, err := svc_acc.NewFromEnv()
	if err != nil {
		panic(err)
	}

	apgo := do.MustInvoke[*appAcc.PostgresOrderer](accounting.Deps())
	do.ProvideValue(injector, apgo)

	// do.ProvideValue(injector, accounting.Config())
	go accounting.Start()

	logger.Slog{}.Info("Test init() - done! Services started")

	return injector
}

func TestOrderThroughKafka(t *testing.T) {
	injector := startServices()

	time.Sleep(20 * time.Millisecond)

	pizClientConf := clients.Config{Host: "localhost", Port: svc_piz.MustEnvConfig().Port}
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

	time.Sleep(20 * time.Millisecond)

	ac := accountingClient.New(clients.Config{Host: "localhost", Port: svc_acc.MustEnvConfig().Port})

	accountingOrder, err := ac.ReadOrder(context.TODO(), createdOrder.ID)
	require.NoError(t, err)

	require.Equal(t, createdOrder.ID, accountingOrder.ID)
	require.Equal(t, accountingOrder.Cost, 5)

	shippingDB := do.MustInvoke[appShipping.OrdererRepository](injector)
	so, err := shippingDB.ReadOrderShippingByID(context.TODO(), createdOrder.ID)
	require.NoError(t, err)
	require.NotEmpty(t, so.ShippedAt)
}
