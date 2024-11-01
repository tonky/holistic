package pizzeria

import (
	"tonky/holistic/infra/kafka"
	"tonky/holistic/infra/logger"

	"github.com/samber/do/v2"
)

func NewDOKafkaFoodOrderProducer(deps do.Injector) (*KafkaFoodOrderProducer, error) {
	config := do.MustInvoke[*kafka.Config](deps)
	logger := do.MustInvoke[*logger.Slog](deps)

	return NewKafkaFoodOrderProducer(*logger, *config)
}
