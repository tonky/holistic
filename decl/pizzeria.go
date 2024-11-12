package decl

import (
	desriber "tonky/holistic/describer"
)

func PizzeriaService() desriber.Service {
	// rpc: net/rpc, twirp, grpc with optional gateway
	getOrder := desriber.Endpoint{
		Name:   "order",
		Method: desriber.Read,
		In:     desriber.Inputs{Name: "food.OrderID"},
		Out: map[desriber.ResponseType]desriber.ResponseObject{
			desriber.ResponseOK: "food.Order",
		},
	}

	createOrder := desriber.Endpoint{
		Name:   "order",
		Method: desriber.Create,
		In:     desriber.Inputs{Name: "NewOrder"},
		Out: map[desriber.ResponseType]desriber.ResponseObject{
			desriber.ResponseOK: "food.Order",
		},
	}

	updateOrder := desriber.Endpoint{
		Name:   "order",
		Method: desriber.Update,
		In:     desriber.Inputs{Name: "UpdateOrder"},
		Out: map[desriber.ResponseType]desriber.ResponseObject{
			desriber.ResponseOK: "food.Order",
		},
	}

	return desriber.Service{
		Name:           "pizzeria",
		Rpc:            desriber.GoNative,
		Dependencies:   desriber.SamberDO,
		Endpoints:      []desriber.Endpoint{getOrder, createOrder, updateOrder},
		ConfigItems:    []desriber.ConfigItem{{Name: "StartDelayMs", Typ: "int"}},
		AppConfigItems: []desriber.ConfigItem{{Name: "EnabledInRegions", Typ: "[]string"}},
		Postgres: []desriber.Postgres{{
			Name: "orderer",
			Methods: []desriber.InterfaceMethod{
				{
					Name: "ReadOrderByID",
					Arg:  desriber.InfraObject{Name: "orderID", Typ: "food.OrderID"},
					Ret:  desriber.InfraObject{Name: "order", Typ: "food.Order"},
				},
				{
					Name: "SaveOrder",
					Arg:  desriber.InfraObject{Name: "newOrder", Typ: "NewOrder"},
					Ret:  desriber.InfraObject{Name: "order", Typ: "food.Order"},
				},
				{
					Name: "UpdateOrder",
					Arg:  desriber.InfraObject{Name: "newOrder", Typ: "UpdateOrder"},
					Ret:  desriber.InfraObject{Name: "order", Typ: "food.Order"},
				},
			},
		}},
		KafkaProducers: []desriber.TopicDesc{TopicFoodOrderCreated, TopicFoodOrderUpdated},
	}
}
