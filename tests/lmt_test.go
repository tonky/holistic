package tests

import (
	"context"
	"os"
	"testing"
	"time"
	"tonky/holistic/domain/food"

	// domain "tonky/holistic/domain/shipping"
	"tonky/holistic/infra/kafkaConsumer"
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/services/shipping"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestServerWithLMT(t *testing.T) {
	os.Setenv("SHIPPING_PORT", "1237")
	os.Setenv("SHIPPING_APP_FOODORDERER_PORT", "5433")
	os.Setenv("SHIPPING_APP_KAFKA_GROUP_ID", "test-order-through-kafka-shipping")

	s, err := shipping.NewFromEnv()
	require.NoError(t, err)

	require.NotNil(t, s.App().LMT)
	require.NotNil(t, s.App().LMT.Logger)
	require.NotNil(t, s.App().LMT.Metrics)
	require.NotNil(t, s.App().LMT.Tracer)

	go s.Start()

	time.Sleep(20 * time.Millisecond)

	op, err := kafkaProducer.NewOrderPaidProducer(s.Deps().LMT, s.Config().App.Kafka)
	require.NoError(t, err)

	osc, err := kafkaConsumer.NewOrderShippedConsumer(s.Deps().LMT, s.Config().App.Kafka)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	shippedOrdersChan, _ := osc.Chan(ctx)

	oid, err := food.NewOrderID(uuid.New().String())
	require.NoError(t, err)

	require.NoError(t, op.ProduceOrderPaid(context.TODO(), food.Order{OrderID: oid, Content: "test", IsFinal: true}))

	shippedOrder := <-shippedOrdersChan
	require.Equal(t, shippedOrder.OrderID, oid)
	require.Equal(t, shippedOrder.Price, 5)
}
