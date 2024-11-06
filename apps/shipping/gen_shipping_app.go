// AUTOGENERATED! DO NOT EDIT.

package shipping

import (
	"context"
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/infra/kafkaConsumer"
	"tonky/holistic/infra/logger"
)

type Deps struct {
	Config Config
	Logger *logger.Slog
    OrdererRepo OrdererRepository
    ShippingOrderShippedProducer kafkaProducer.IShippingOrderShipped
    AccountingOrderPaidConsumer kafkaConsumer.IAccountingOrderPaid
}


type App struct {
	Deps		Deps
	Logger		*logger.Slog
}

func NewApp(deps Deps) (App, error) {
	if deps.Logger == nil {
		deps.Logger = &logger.Slog{}
	}

	ctx := context.Background()

	app := App{
		Deps:       deps,
		Logger:     deps.Logger,
	}

	go func() {
		for err := range app.Deps.AccountingOrderPaidConsumer.Run(ctx, app.AccountingOrderPaidProcessor) {
			app.Deps.Logger.Warn(err.Error())
		}
	}()
	
	return app, nil
}
