package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"tonky/holistic/clients/pizzeriaClient"
	"tonky/holistic/domain/food"
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/kafkaConsumer"
	"tonky/holistic/infra/logger"
	"tonky/holistic/services/accounting"
	"tonky/holistic/services/pizzeria"
	svc_piz "tonky/holistic/services/pizzeria"

	"github.com/stretchr/testify/require"
)

// order
// created -> finalized -> paid -> shipped -> delivered

func TestE2E(t *testing.T) {
	os.Setenv("PIZZERIA_PORT", "1238")
	os.Setenv("PIZZERIA_APP_POSTGRESORDERER_PORT", "5432")
	os.Setenv("PIZZERIA_APP_KAFKA_GROUP_ID", "test-e2e-pizzeria")

	os.Setenv("ACCOUNTING_PORT", "1239")
	os.Setenv("ACCOUNTING_APP_POSTGRESORDERER_PORT", "5434")
	os.Setenv("ACCOUNTING_APP_KAFKA_GROUP_ID", "test-e2e-accounting")

	os.Setenv("SHIPPING_APP_KAFKA_GROUP_ID", "test-order-through-kafka-shipping")

	svcPiz, err := pizzeria.NewFromEnv()
	require.NoError(t, err)

	svcAcc, err := accounting.NewFromEnv()
	require.NoError(t, err)

	go svcPiz.Start()
	go svcAcc.Start()

	time.Sleep(30 * time.Millisecond)

	pc := pizzeriaClient.NewFromEnv("dev")

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

	kConf := kafka.Config{
		Brokers: svcPiz.Config().App.Kafka.Brokers,
		GroupID: "test-e2e-test",
	}

	consumer, err := kafkaConsumer.NewFoodOrderUpdatedConsumer(logger.Slog{}, kConf)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	processor := func(ctx context.Context, order food.Order) error {
		fmt.Printf("processor got food order: %+v\n", order)
		require.Equal(t, order.ID, createdOrder.ID)

		return nil
	}

	errCh := consumer.Run(ctx, processor)

	updatedOrder, err := pc.UpdateOrder(context.TODO(), updateOrder)
	require.NoError(t, err)
	require.EqualValues(t, updateOrder, updatedOrder)

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-ctx.Done():
	}

	// require.True(t, false)
}
