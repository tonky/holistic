// AUTOGENERATED! DO NOT EDIT.
package shippingClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"tonky/holistic/clients"
	 "tonky/holistic/domain/food"
	 "tonky/holistic/domain/shipping"
)

type IShippingClient interface {
	ReadOrder(context.Context, food.OrderID) (shipping.Order, error)
}

func New(config clients.Config) ShippingClient {
	return ShippingClient{
		config: config,
	}
}

func NewFromEnv(env string) ShippingClient {
	envConf := clients.ConfigForEnv("shipping", env)

	return ShippingClient{
		config: envConf,
	}
}

type ShippingClient struct {
	config clients.Config
}

func (c ShippingClient) ReadOrder(ctx context.Context, arg food.OrderID) (shipping.Order, error) {
	var reply shipping.Order

	jsonBody, err := json.Marshal(arg)
	if err != nil { return reply, err}

	bodyReader := bytes.NewReader(jsonBody)

 	requestURL := fmt.Sprintf("%s/%s", c.config.ServerAddress(), "GetOrder")

 	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil { return reply, err }

	res, err := http.DefaultClient.Do(req)
	if err != nil { return reply, err }

	resBody, err := io.ReadAll(res.Body)
	if err != nil { return reply, err }

	if err := json.Unmarshal(resBody, &reply); err != nil { return reply, err }

	return reply, nil
}

