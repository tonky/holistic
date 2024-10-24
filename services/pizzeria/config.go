// AUTOGENERATED! DO NOT EDIT.
package pizzeria

import (
	"tonky/holistic/infra/logger"

	"tonky/holistic/infra/postgres"

	"github.com/kelseyhightower/envconfig"
	"github.com/samber/do/v2"
)

// service specific config - env, secrets, run mode, flags etc
type Config struct {
	Environment   string `default:"dev"`
	Port int `default:"1234"`

    Postgres postgres.Config

    ShouldMockApp bool `split_words:"true"`
}

func NewEnvConfig() (Config, error) {
	var c Config

	return c, envconfig.Process("pizzeria", &c)
}

func NewConfig(i do.Injector) (*Config, error) {
	config, err := NewEnvConfig()
	if err != nil {
		return nil, err
	}

	do.ProvideValue(i, &config.Postgres)
	do.Provide(i, logger.NewSlogLogger)

	return &config, nil
}
