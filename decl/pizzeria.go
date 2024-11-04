package decl

import (
	"tonky/holistic/generator/services"
)

func PizzeriaService() services.Service {
	// rpc: net/rpc, twirp, grpc with optional gateway
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

	updateOrder := services.Endpoint{
		Name:   "order",
		Method: services.Update,
		In:     services.Inputs{Name: "UpdateOrder"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK: "food.Order",
		},
	}

	return services.Service{
		Name:           "pizzeria",
		Rpc:            services.GoNative,
		Endpoints:      []services.Endpoint{getOrder, createOrder, updateOrder},
		ConfigItems:    []services.ConfigItem{{Name: "StartDelayMs", Typ: "int"}},
		AppConfigItems: []services.ConfigItem{{Name: "EnabledInRegions", Typ: "[]string"}},
		Postgres: []services.Postgres{{
			Name: "orderer",
			Methods: []services.InterfaceMethod{
				{
					Name: "ReadOrderByID",
					Arg:  services.InfraObject{Name: "orderID", Typ: "food.OrderID"},
					Ret:  services.InfraObject{Name: "order", Typ: "food.Order"},
				},
				{
					Name: "SaveOrder",
					Arg:  services.InfraObject{Name: "newOrder", Typ: "NewOrder"},
					Ret:  services.InfraObject{Name: "order", Typ: "food.Order"},
				},
			},
		}},
		KafkaProducers: []services.TopicDesc{services.TopicFoodOrderCreated},
	}
}
