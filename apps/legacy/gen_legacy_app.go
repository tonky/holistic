// AUTOGENERATED! DO NOT EDIT.

package legacy

import (
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/clients/accountingClient"
	"tonky/holistic/infra/logger"
)

type Deps struct {
	Config Config
	Logger *logger.Slog
    OrdererRepo OrdererRepository
    FoodOrderCreatedProducer kafkaProducer.IFoodOrderCreated
    FoodOrderUpdatedProducer kafkaProducer.IFoodOrderUpdated
}

type Clients struct {
    AccountingClient accountingClient.IAccountingClient
}

type App struct {
	Deps		Deps
	Clients		Clients
	Logger		*logger.Slog
}

func NewApp(deps Deps, clients Clients) (App, error) {
	if deps.Logger == nil {
		deps.Logger = &logger.Slog{}
	}

	app := App{
		Deps:       deps,
		Logger:     deps.Logger,
		Clients: 	clients,
	}

	return app, nil
}
