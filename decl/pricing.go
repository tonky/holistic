package decl

import (
	describer "tonky/holistic/describer"
)

func PricingService() describer.Service {
	getOrderPrice := describer.Endpoint{
		Name:   "order",
		Method: describer.Read,
		In:     describer.Inputs{Name: "food.OrderID"},
		Out: map[describer.ResponseType]describer.ResponseObject{
			describer.ResponseOK: "pricing.OrderPrice",
		},
	}

	return describer.Service{
		Name:         "pricing",
		Rpc:          describer.HTTP,
		Dependencies: describer.Struct,
		Endpoints:    []describer.Endpoint{getOrderPrice},
		Postgres: []describer.Postgres{{
			Name: "orderer",
			Methods: []describer.InterfaceMethod{
				{
					Name: "ReadOrderPricingByID",
					Arg:  describer.InfraObject{Name: "orderID", Typ: "food.OrderID"},
					Ret:  describer.InfraObject{Name: "pricing", Typ: "pricing.OrderPrice"},
				},
			},
		}},
	}
}
