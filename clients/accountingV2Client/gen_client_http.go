// AUTOGENERATED! DO NOT EDIT.
package accountingV2Client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/slogLogger"
	"tonky/holistic/domain/foodStore"
	"tonky/holistic/domain/accountingV2"
	"tonky/holistic/clients"
	svc "tonky/holistic/services/accountingV2"
)

type IAccountingV2Client interface {
	GetOrderByID(context.Context, foodStore.OrderID) (foodStore.Order, error)
	CreateOrder(context.Context, accountingV2.NewFoodOrder) (foodStore.Order, error)
}

func New(config clients.Config) AccountingV2Client {
	return AccountingV2Client {
		config: config,
		logger: slogLogger.Default(),
	}
}

func NewFromEnv() AccountingV2Client {
	svcConf := svc.MustEnvConfig()

	return AccountingV2Client {
		config: clients.Config{Host: "http://localhost", Port: svcConf.Port},
        logger: slogLogger.NewFromConfig(svcConf.Logger),
	}
}

type AccountingV2Client struct {
	config clients.Config
    logger logger.ILogger
}

func (c AccountingV2Client) GetOrderByID(ctx context.Context, arg foodStore.OrderID) (foodStore.Order, error) {
	var reply foodStore.Order

	jsonBody, err := json.Marshal(arg)
	if err != nil { return reply, err}

	bodyReader := bytes.NewReader(jsonBody)

 	requestURL := fmt.Sprintf("%s/%s", c.config.ServerAddress(), "GetOrderByID")

	c.logger.Debug("AccountingV2Client.GetOrderByID()", "requestURL", requestURL, "newRecipe", arg)

 	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil { return reply, err }

	res, err := http.DefaultClient.Do(req)
	if err != nil { return reply, err }

	resBody, err := io.ReadAll(res.Body)
	if err != nil { return reply, err }

	if err := json.Unmarshal(resBody, &reply); err != nil { return reply, err }

	return reply, nil
}

func (c AccountingV2Client) CreateOrder(ctx context.Context, arg accountingV2.NewFoodOrder) (foodStore.Order, error) {
	var reply foodStore.Order

	jsonBody, err := json.Marshal(arg)
	if err != nil { return reply, err}

	bodyReader := bytes.NewReader(jsonBody)

 	requestURL := fmt.Sprintf("%s/%s", c.config.ServerAddress(), "CreateOrder")

	c.logger.Debug("AccountingV2Client.CreateOrder()", "requestURL", requestURL, "newRecipe", arg)

 	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil { return reply, err }

	res, err := http.DefaultClient.Do(req)
	if err != nil { return reply, err }

	resBody, err := io.ReadAll(res.Body)
	if err != nil { return reply, err }

	if err := json.Unmarshal(resBody, &reply); err != nil { return reply, err }

	return reply, nil
}

