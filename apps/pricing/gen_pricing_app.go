// AUTOGENERATED! DO NOT EDIT.

package pricing

import (
	"tonky/holistic/infra/logger"
)

type Deps struct {
	Config Config
	Logger logger.ILogger
    OrdererRepo OrdererRepository
}


type App struct {
	Deps		Deps
}

func NewApp(deps Deps) (App, error) {
	app := App{
		Deps:       deps,
	}

	return app, nil
}

