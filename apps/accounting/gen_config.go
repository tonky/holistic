// AUTOGENERATED! DO NOT EDIT.
package accounting

import (
	"tonky/holistic/infra/kafka"

	"github.com/kelseyhightower/envconfig"
)

// app specific config - infra, flags etc
type Config struct {
	Environment   string `default:"dev"`

    Kafka kafka.Config

}

func NewEnvConfig() (Config, error) {
	var c Config

	return c, envconfig.Process("accounting", &c)
}
