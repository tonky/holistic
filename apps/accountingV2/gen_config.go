// AUTOGENERATED! DO NOT EDIT.
package accountingV2

import (
	"tonky/holistic/infra/kafka"

	"github.com/kelseyhightower/envconfig"
)

// app specific config - infra, flags etc
type Config struct {
	Environment   string `default:"dev"`

    Kafka kafka.Config

    KafkaConsumptionRPSLimit int `split_words:"false"`
}

func NewEnvConfig() (Config, error) {
	var c Config

	return c, envconfig.Process("accountingV2", &c)
}