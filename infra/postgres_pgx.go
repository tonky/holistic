package infra

import (
	"tonky/holistic/infra/postgres"

	"github.com/samber/do/v2"
)

func NewPostgresClient(i do.Injector) (*postgres.Client, error) {
	conf := do.MustInvoke[*postgres.Config](i)

	pgc, err := postgres.NewClient(*conf)

	return &pgc, err
}
