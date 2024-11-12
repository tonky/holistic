package decl

import (
	"tonky/holistic/describer"
)

func ShippingService() describer.Service {
	getShipment := describer.Endpoint{
		Name:   "order",
		Method: describer.Read,
		In:     describer.Inputs{Name: "food.OrderID"},
		Out: map[describer.ResponseType]describer.ResponseObject{
			describer.ResponseOK: "shipping.Order",
		},
	}

	return describer.Service{
		Name:           "shipping",
		Rpc:            describer.HTTP,
		Dependencies:   describer.Struct,
		Endpoints:      []describer.Endpoint{getShipment},
		KafkaProducers: []describer.TopicDesc{TopicShippingOrderShipped},
		KafkaConsumers: []describer.TopicDesc{TopicAccountingOrderPaid},
		Postgres: []describer.Postgres{{
			Name: "orderer",
			Methods: []describer.InterfaceMethod{
				{
					Name: "ReadOrderShippingByID",
					Arg:  describer.InfraObject{Name: "orderID", Typ: "food.OrderID"},
					Ret:  describer.InfraObject{Name: "order", Typ: "shipping.Order"},
				},
				{
					Name: "SaveShipping",
					Arg:  describer.InfraObject{Name: "shipping", Typ: "shipping.Order"},
					Ret:  describer.InfraObject{Name: "order", Typ: "shipping.Order"},
				},
			},
		}},
	}
}
