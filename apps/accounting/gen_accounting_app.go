// AUTOGENERATED! DO NOT EDIT.

package accounting

import (
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/infra/kafkaConsumer"
	"context"
	"tonky/holistic/clients/pricingClient"
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

type App struct {
	Deps       do.Injector
	Logger     *logger.Slog

    OrdererRepo OrdererRepository
    AccountingOrderPaidProducer kafkaProducer.IAccountingOrderPaid
    FoodOrderUpdatedConsumer kafkaConsumer.IFoodOrderUpdated
    PricingClient pricingClient.IPricingClient
}

func NewApp(deps do.Injector) (*App, error) {
	app := App{
		Deps:       deps,
		Logger:     do.MustInvoke[*logger.Slog](deps),
        OrdererRepo: do.MustInvokeAs[OrdererRepository](deps),
        AccountingOrderPaidProducer: do.MustInvokeAs[kafkaProducer.IAccountingOrderPaid](deps),
        FoodOrderUpdatedConsumer: do.MustInvokeAs[kafkaConsumer.IFoodOrderUpdated](deps),
        PricingClient: do.MustInvokeAs[pricingClient.IPricingClient](deps),

	}

	return &app, nil
}

func (a App) RunConsumers() {
	a.Logger.Info(">> accounting.App.RunConsumers()")

	ctx := context.Background()

	go func() {
		for err := range a.FoodOrderUpdatedConsumer.Run(ctx, a.FoodOrderUpdatedProcessor) {
			a.Logger.Warn(err.Error())
		}
	}()
}
