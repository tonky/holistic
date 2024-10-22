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
		In:     services.Inputs{Name: "food.Order"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK: "food.Order",
		},
	}

	return services.Service{
		Name:        "pizzeria",
		Rpc:         services.GoNative,
		Endpoints:   []services.Endpoint{getOrder, createOrder},
		ConfigItems: []services.ConfigItem{{Name: "ShouldMockApp", Typ: "bool"}},
		Infra: []services.Infra{
			{
				Name: "Orderer",
				Typ:  "postgres",
				InOut: []services.InOut{
					{
						Name: "ReadOrderByID",
						In:   services.InfraObject{Name: "orderID", Typ: "food.OrderID"},
						Out:  services.InfraObject{Name: "order", Typ: "food.Order"},
					},
					{
						Name: "SaveOrder",
						In:   services.InfraObject{Name: "newOrder", Typ: "NewOrder"},
						Out:  services.InfraObject{Name: "order", Typ: "food.Order"},
					},
				},
			},
		},
	}
}
