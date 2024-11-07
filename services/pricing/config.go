// AUTOGENERATED! DO NOT EDIT.
package pricing

import (
	app "tonky/holistic/apps/pricing"
	"github.com/kelseyhightower/envconfig"
	"github.com/samber/do/v2"
)

// service specific config - env, secrets, run mode, flags etc
type Config struct {
	Environment   string `default:"dev"`
	Port int `default:"1234"`

	App app.Config

}

func NewEnvConfig() (Config, error) {
	var c Config

	return c, envconfig.Process("pricing", &c)
}

func MustEnvConfig() Config {
	conf, err := NewEnvConfig()
	if err != nil { panic(err) }

	return conf
}

func NewConfig(i do.Injector) (*Config, error) {
	config, err := NewEnvConfig()
	if err != nil {
		return nil, err
	}

	return &config, nil
}
