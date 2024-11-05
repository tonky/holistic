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

type Clients struct {
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
