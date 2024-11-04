package main

import (
	"fmt"
	"net/http"
	app "tonky/holistic/apps/legacy"
	"tonky/holistic/clients/accountingClient"
	"tonky/holistic/infra/logger"
	svc "tonky/holistic/services/legacy"
)

func main() {
	l := logger.Slog{}

	conf, err := svc.NewEnvConfig()
	if err != nil {
		panic(err)
	}

	deps := app.NewDeps(conf.App)

	pg, err := app.NewPostgresOrdererRepository(l, conf.App.PostgresOrderer)
	if err != nil {
		panic(err)
	}

	deps.OrdererRepo = pg

	clients := app.Clients{
		AccountingClient: accountingClient.NewMock(),
	}

	routes, err := svc.NewLegacy(conf, deps, clients)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server on port", conf.Port)

	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), routes)
}
