package main

import (
	app "tonky/holistic/apps/pizzeria"
	"tonky/holistic/infra"
	"tonky/holistic/services/pizzeria"

	"github.com/samber/do/v2"
)

func main() {
	injector := do.New()
	do.Provide(injector, pizzeria.NewConfig)

	// provide infra dependencies
	do.Provide(injector, infra.NewPostgresClient)

	do.Provide(injector, app.NewPostgresOrdererRepository)

	svc := pizzeria.NewPizzeria(injector)

	svc.Start()
}
