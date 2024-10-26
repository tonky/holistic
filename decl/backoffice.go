package decl

import "tonky/holistic/generator/services"

func BackofficeService() services.Service {
	getOrder := services.Endpoint{
		Name:   "order",
		Method: services.Read,
		In:     services.Inputs{Name: "food.OrderID"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK: "accounting.Order",
		},
	}
	return services.Service{
		Name:        "accounting",
		Rpc:         services.GoNative,
		ConfigItems: []services.ConfigItem{{Name: "KafkaConsumptionRPS", Typ: "int"}},
		Endpoints:   []services.Endpoint{getOrder},
		KafkaConsumers: []services.KafkaConsumer{{
			Name:  "foodOrder",
			Topic: "pizzeria.order",
			Model: "food.Order",
		}},
	}
}
