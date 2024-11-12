// AUTOGENERATED! DO NOT EDIT.

package pizzeria

import (
	"tonky/holistic/infra/kafkaProducer"
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

type Deps do.Injector

type App struct {
	Deps       Deps
	Logger     *logger.Slog
    OrdererRepo OrdererRepository
    FoodOrderCreatedProducer kafkaProducer.IFoodOrderCreated
    FoodOrderUpdatedProducer kafkaProducer.IFoodOrderUpdated
}

func NewApp(deps Deps) (*App, error) {
	app := App{
		Deps:       deps,
		Logger:     do.MustInvoke[*logger.Slog](deps),
        OrdererRepo: do.MustInvokeAs[OrdererRepository](deps),
        FoodOrderCreatedProducer: do.MustInvokeAs[kafkaProducer.IFoodOrderCreated](deps),
        FoodOrderUpdatedProducer: do.MustInvokeAs[kafkaProducer.IFoodOrderUpdated](deps),
	}

	return &app, nil
}

