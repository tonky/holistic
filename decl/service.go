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
		Infra:       []services.Infra{{Name: "postgres"}},
	}
}
