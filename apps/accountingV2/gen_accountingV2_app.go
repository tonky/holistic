// AUTOGENERATED! DO NOT EDIT.

package accountingV2

import (
	"tonky/holistic/infra/logger"
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/infra/kafkaConsumer"
	"context"
	"tonky/holistic/clients/pricingClient"
)

type Clients struct {
    PricingClient pricingClient.IPricingClient
}

type Deps struct {
	Config Config
	Logger *logger.Slog
    AccountingOrderPaidProducer kafkaProducer.IAccountingOrderPaid
    FoodOrderUpdatedConsumer kafkaConsumer.IFoodOrderUpdated
}

type App struct {
	Deps       Deps
	Logger     *logger.Slog
	Clients		Clients
}

func NewApp(deps Deps, clients Clients) (*App, error) {
	app := App{
		Deps:       deps,
		Clients: clients,
		Logger:     deps.Logger,
	}

	return &app, nil
}

func (a App) RunConsumers() {
	a.Logger.Info(">> accountingV2.App.RunConsumers()")

	ctx := context.Background()

	go func() {
		for err := range a.Deps.FoodOrderUpdatedConsumer.Run(ctx, a.FoodOrderUpdatedProcessor) {
			a.Logger.Warn(err.Error())
		}
	}()
}