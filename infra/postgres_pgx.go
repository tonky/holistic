package infra

import (
	"tonky/holistic/infra/postgres"
)

func NewPostgresClient(conf postgres.Config) (postgres.Client, error) {
	return postgres.NewClient(conf)
}
