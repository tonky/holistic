package main

import (
	"tonky/holistic/services/pizzeria"

	"github.com/samber/do/v2"
)

func main() {
	injector := do.New()
	do.Provide(injector, pizzeria.NewConfig)

	// provide infra dependencies
	// do.Provide(injector, NewPostgresClient)

	svc := pizzeria.NewPizzeria(injector)

	svc.Start()
}

/*
func NewPostgresClient(i do.Injector) (*infra.PostgresClient, error) {
	conf := do.MustInvoke[*pizzeria.Config](i)

	pgc, err := infra.NewPostgresClient(conf.Postgres)

	return &pgc, err
}
*/
