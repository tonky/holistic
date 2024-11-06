// AUTOGENERATED! DO NOT EDIT.

package pricing

import (
	"tonky/holistic/infra/logger"
)

type Deps struct {
	Config Config
	Logger *logger.Slog
    OrdererRepo OrdererRepository
}


type App struct {
	Deps		Deps
	Logger		*logger.Slog
}

func NewApp(deps Deps) (App, error) {
	if deps.Logger == nil {
		deps.Logger = &logger.Slog{}
	}

	app := App{
		Deps:       deps,
		Logger:     deps.Logger,
	}

	return app, nil
}
