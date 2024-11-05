package decl

import (
	"tonky/holistic/generator/services"
)

func ShippingService() services.Service {
	return services.Service{
		Name: "shippng",
		// Rpc:            services.HTTP,
		KafkaProducers: []services.TopicDesc{services.TopicAccountingOrderPaid},
		KafkaConsumers: []services.TopicDesc{services.TopicShippingOrderShipped},
	}
}
