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
	"tonky/holistic/infra/logger"

	svc "tonky/holistic/services/legacy"

	"github.com/stretchr/testify/require"
)

func TestUpdateOrder(t *testing.T) {
	memProducer, kafkaMemStore := app.NewMemoryOrderProducer()

	appDeps := app.Deps{
		OrdererRepo:       app.NewMemoryOrdererRepository(),
		FoodOrderProducer: memProducer,
		Logger:            &logger.Slog{},
	}

	acm := accountingClient.NewMock()

	appClients := app.Clients{
		AccountingClient: acm,
	}

	conf := svc.Config{Port: 1237}

	go runServer(conf, appDeps, appClients)

	time.Sleep(100 * time.Millisecond)

	c := legacyClient.New(clients.Config{Host: "http://localhost", Port: conf.Port})

	createdOrder, err := c.CreateOrder(context.TODO(), svc.NewOrder{Content: "legacy_content"})
	require.NoError(t, err)
	require.Equal(t, "legacy_content", createdOrder.Content)

	require.Len(t, *kafkaMemStore, 0, kafkaMemStore)

	// setup accounting mock for client call response
	acm.Orders[createdOrder.ID] = accounting.Order{ID: createdOrder.ID, Content: "legacy_content", Cost: 4, IsPaid: true}

	_, errC := c.UpdateOrder(context.TODO(), svc.UpdateOrder{ID: createdOrder.ID.String(), Content: "legacy_content update", IsFinal: true})
	require.NoError(t, errC)

	require.Len(t, *kafkaMemStore, 1, kafkaMemStore)
	require.True(t, false)
}

func runServer(conf svc.Config, appDeps app.Deps, clients app.Clients) error {
	routes, err := svc.NewLegacy(conf, appDeps, clients)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server on port", conf.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), routes)
}
