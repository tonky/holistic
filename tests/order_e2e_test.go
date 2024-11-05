package tests

import (
	"context"
	"os"
	"testing"
	"time"
	"tonky/holistic/clients"
	"tonky/holistic/clients/pizzeriaClient"
	"tonky/holistic/services/pizzeria"
	svc_piz "tonky/holistic/services/pizzeria"

	"github.com/stretchr/testify/require"
)

// order
// created -> finalized -> paid -> shipped -> delivered

func TestOrderCreated(t *testing.T) {
	os.Setenv("PIZZERIA_PORT", "1237")
	svc, err := pizzeria.NewFromEnv()
	if err != nil {
		panic(err)
	}

	go svc.Start()

	time.Sleep(30 * time.Millisecond)

	pc := pizzeriaClient.New(clients.Config{Port: svc.Config().Port})

	newOrder := svc_piz.NewOrder{
		Content: "new order",
	}

	createdOrder, err := pc.CreateOrder(context.TODO(), newOrder)
	require.NoError(t, err)
	require.Equal(t, newOrder.Content, createdOrder.Content)
}
