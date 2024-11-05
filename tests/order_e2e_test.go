package tests

import (
	"context"
	"os"
	"testing"
	"time"
	"tonky/holistic/clients"
	"tonky/holistic/clients/pizzeriaClient"
	"tonky/holistic/services/accounting"
	"tonky/holistic/services/pizzeria"
	svc_piz "tonky/holistic/services/pizzeria"

	"github.com/stretchr/testify/require"
)

// order
// created -> finalized -> paid -> shipped -> delivered

func TestE2E(t *testing.T) {
	os.Setenv("PIZZERIA_PORT", "1237")
	os.Setenv("ACCOUNTING_PORT", "1238")

	svcPiz, err := pizzeria.NewFromEnv()
	require.NoError(t, err)

	svcAcc, err := accounting.NewFromEnv()
	require.NoError(t, err)

	go svcPiz.Start()
	go svcAcc.Start()

	time.Sleep(30 * time.Millisecond)

	pc := pizzeriaClient.New(clients.Config{Port: svcPiz.Config().Port})

	newOrder := svc_piz.NewOrder{
		Content: "new order",
	}

	createdOrder, err := pc.CreateOrder(context.TODO(), newOrder)
	require.NoError(t, err)
	require.Equal(t, newOrder.Content, createdOrder.Content)

	updateOrder := svc_piz.UpdateOrder{
		ID:      createdOrder.ID,
		Content: "updated content",
		IsFinal: true,
	}

	updatedOrder, err := pc.UpdateOrder(context.TODO(), updateOrder)
	require.NoError(t, err)
	require.EqualValues(t, updateOrder, updatedOrder)
}
