package decl

import (
	"tonky/holistic/generator/services"
)

func ShippingService() services.Service {
	getShipment := services.Endpoint{
		Name:   "order",
		Method: services.Read,
		In:     services.Inputs{Name: "food.OrderID"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK: "shipping.Order",
		},
	}

	return services.Service{
		Name:           "shipping",
		Rpc:            services.HTTP,
		Endpoints:      []services.Endpoint{getShipment},
		KafkaProducers: []services.TopicDesc{services.TopicShippingOrderShipped},
		KafkaConsumers: []services.TopicDesc{services.TopicAccountingOrderPaid},
		Postgres: []services.Postgres{{
			Name: "orderer",
			Methods: []services.InterfaceMethod{
				{
					Name: "ReadOrderShippingByID",
					Arg:  services.InfraObject{Name: "orderID", Typ: "food.OrderID"},
					Ret:  services.InfraObject{Name: "order", Typ: "shipping.Order"},
				},
				{
					Name: "SaveShipping",
					Arg:  services.InfraObject{Name: "shipping", Typ: "shipping.Order"},
					Ret:  services.InfraObject{Name: "order", Typ: "shipping.Order"},
				},
			},
		}},
	}
}
