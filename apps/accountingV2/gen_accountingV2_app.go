// AUTOGENERATED! DO NOT EDIT.

package accountingV2

import (
	"log/slog"
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
    AccountingOrderPaidProducer kafkaProducer.IAccountingOrderPaid
    FoodOrderUpdatedConsumer kafkaConsumer.IFoodOrderUpdated
    FoodOrderer IFoodOrderer
	Clients		Clients
}

type App struct {
	Deps       Deps
}

func NewApp(deps Deps) (*App, error) {
	app := App{
		Deps:       deps,
	}

	return &app, nil
}

func MustDepsFromEnv() Deps {
    slog.Debug("accountingV2.App.MustDepsFromenv()")

	cfg := MustEnvConfig()

	deps, err := DepsFromConf(cfg)
	if err != nil {
    	slog.Error("DepsFromConf error", slog.Any("config", cfg), slog.Any("err", err))
	}

	return deps
}

func DepsFromConf(cfg Config) (Deps, error) {
    slog.Debug("accountingV2.App.DepsFromConf()", slog.Any("config", cfg))

    deps := Deps{}

    FoodOrderer, err := NewFoodOrderer(logger.Slog{}, cfg.FoodOrderer)
    if err != nil {
        return deps, err
    }

    deps.FoodOrderer = FoodOrderer

	AccountingOrderPaidProducer, err := kafkaProducer.NewAccountingOrderPaidProducer(logger.Slog{}, cfg.Kafka)
    if err != nil {
        return deps, err
    }
	deps.AccountingOrderPaidProducer = AccountingOrderPaidProducer
	FoodOrderUpdatedConsumer, err := kafkaConsumer.NewFoodOrderUpdatedConsumer(logger.Slog{}, cfg.Kafka)
    if err != nil {
        return deps, err
    }
	deps.FoodOrderUpdatedConsumer = FoodOrderUpdatedConsumer

    deps.Clients = Clients{
        PricingClient: pricingClient.NewFromEnv(cfg.Environment),
    }

    return deps, nil
}

func (a App) RunConsumers() {
	slog.Info("accountingV2.App.RunConsumers()")

	ctx := context.Background()

	go func() {
		for err := range a.Deps.FoodOrderUpdatedConsumer.Run(ctx, a.FoodOrderUpdatedProcessor) {
			slog.Warn(err.Error())
		}
	}()
}
