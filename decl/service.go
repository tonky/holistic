package decl

import (
	"tonky/holistic/generator/services"
)

func PizzeriaService() services.Service {
	// rpc: net/rpc, twirp
	getOrder := services.Endpoint{
		Name:   "order",
		Method: services.Read,
		In:     services.Inputs{Name: "food.OrderID"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK: "food.Order",
		},
	}

	createOrder := services.Endpoint{
		Name:   "order",
		Method: services.Create,
		In:     services.Inputs{Name: "NewOrder"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK: "food.Order",
		},
	}

	return services.Service{
		Name:        "pizzeria",
		Rpc:         services.GoNative,
		Endpoints:   []services.Endpoint{getOrder, createOrder},
		ConfigItems: []services.ConfigItem{{Name: "ShouldMockApp", Typ: "bool"}},
		Postgres: []services.Postgres{{
			Name: "orderer",
			Interface: []services.Interface{
				{
					Method: "ReadOrderByID",
					Arg:    services.InfraObject{Name: "orderID", Typ: "food.OrderID"},
					Ret:    services.InfraObject{Name: "order", Typ: "food.Order"},
				},
				{
					Method: "SaveOrder",
					Arg:    services.InfraObject{Name: "newOrder", Typ: "NewOrder"},
					Ret:    services.InfraObject{Name: "order", Typ: "food.Order"},
				},
			},
		}},
		KafkaProducers: []services.KafkaProducer{{
			Name:  "foodOrder",
			Topic: "pizzeria.order",
			Model: "food.Order",
		}},
		/*
			Infra: []services.Infra{
				{
					Name: "OrderProducer",
					Typ:  "kafka",
					InOut: []services.InOut{
						{
							Name: "ProduceNewOrder",
							Out:  services.InfraObject{Name: "pizzeria.orders", Typ: "food.Order"},
						},
					},
				},
			},
		*/
	}
}
