package decl

import (
	"tonky/holistic/generator/services"
)

func PricingService() services.Service {
	getOrderPrice := services.Endpoint{
		Name:   "order",
		Method: services.Read,
		In:     services.Inputs{Name: "food.OrderID"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK: "pricing.OrderPrice",
		},
	}

	return services.Service{
		Name:      "pricing",
		Rpc:       services.HTTP,
		Endpoints: []services.Endpoint{getOrderPrice},
		Postgres: []services.Postgres{{
			Name: "orderer",
			Methods: []services.InterfaceMethod{
				{
					Name: "ReadOrderPricingByID",
					Arg:  services.InfraObject{Name: "orderID", Typ: "food.OrderID"},
					Ret:  services.InfraObject{Name: "pricing", Typ: "pricing.OrderPrice"},
				},
			},
		}},
	}
}
