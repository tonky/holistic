package {{ service_name }}

import (
	"tonky/holistic/configs/infra"

	"github.com/kelseyhightower/envconfig"
	"github.com/samber/do/v2"
)

// service specific config - env, secrets, run mode, flags etc
type Config struct {
	Environment   string `default:"dev"`
	Port int `default:"1234"`

    {% for i in infra %}
    {{ cap(i.Name) }} infra.{{ i.ConfigVar() }}
    {% end %}

    {% for configItem in config_items %}
    {{ configItem.Name }} {{ configItem.Typ }} `split_words:"true"`
    {% end %}
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

	return &config, nil
}
