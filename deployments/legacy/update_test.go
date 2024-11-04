package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
	app "tonky/holistic/apps/legacy"
	"tonky/holistic/clients"
	"tonky/holistic/clients/accountingClient"
	"tonky/holistic/clients/legacyClient"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/infra/logger"

	svc "tonky/holistic/services/legacy"

	"github.com/stretchr/testify/require"
)

func TestUpdateOrder(t *testing.T) {
	conf := svc.Config{Port: 1237, Environment: "local"}

	memProducer := app.NewMemoryOrderProducer()
	appDeps := app.Deps{
		FoodOrderCreatedProducer: memProducer,
		OrdererRepo:              app.NewMemoryOrdererRepository(),
	}

	kConf := kafka.Config{Brokers: []string{"localhost:19092"}}

	kpou, err := kafkaProducer.NewFoodOrderUpdated(logger.Slog{}, kConf)
	require.NoError(t, err)

	appDeps.FoodOrderUpdatedProducer = kpou

	acm := accountingClient.NewMock()

	appClients := app.Clients{
		AccountingClient: acm,
	}

	go runServer(conf, appDeps, appClients)

	time.Sleep(100 * time.Millisecond)

	c := legacyClient.New(clients.Config{Host: "http://localhost", Port: conf.Port})

	createdOrder, err := c.CreateOrder(context.TODO(), svc.NewOrder{Content: "legacy_content"})
	require.NoError(t, err)
	require.Equal(t, "legacy_content", createdOrder.Content)

	require.Len(t, memProducer.Orders, 0)

	// setup accounting mock for client call response
	acm.Orders[createdOrder.ID] = accounting.Order{ID: createdOrder.ID, Content: "legacy_content", Cost: 4, IsPaid: true}

	_, errC := c.UpdateOrder(context.TODO(), svc.UpdateOrder{ID: createdOrder.ID.String(), Content: "legacy_content update", IsFinal: true})
	require.NoError(t, errC)

	require.Len(t, memProducer.Orders, 1)
	// require.True(t, false)
}

func runServer(conf svc.Config, appDeps app.Deps, clients app.Clients) error {
	routes, err := svc.NewLegacy(conf, appDeps, clients)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server on port", conf.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), routes)
}
