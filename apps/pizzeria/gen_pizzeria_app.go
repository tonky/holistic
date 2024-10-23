// AUTOGENERATED! DO NOT EDIT.

package pizzeria

import (
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

type App struct {
	deps       do.Injector
	logger     *logger.SlogLogger

    ordererRepo OrdererRepository
}

func NewApp(deps do.Injector) (*App, error) {
	return &App{
		deps:       deps,
		logger:     do.MustInvoke[*logger.SlogLogger](deps),
        ordererRepo: do.MustInvokeAs[OrdererRepository](deps),
	}, nil
}
