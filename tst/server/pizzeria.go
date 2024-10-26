package main

import (
	"tonky/holistic/services/pizzeria"

	"github.com/samber/do/v2"
)

func main() {
	injector := do.New()
	do.Provide(injector, pizzeria.NewConfig)

	// provide infra dependencies
	// do.Provide(injector, infra.NewPostgresClient)

	// do.Provide(injector, app.NewPostgresOrdererRepository)

	svc, err := pizzeria.NewPizzeria(injector)
	if err != nil {
		panic(err)
	}

	svc.Start()
}
