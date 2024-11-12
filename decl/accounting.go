package decl

import (
	"tonky/holistic/describer"
)

type Model struct {
	Scope    string
	Domain   string
	Name     string
	External struct {
		PackagePath  string
		PackageAlias string
	}
}

func AccountingService() describer.Service {
	getOrder := describer.Endpoint{
		Name:   "order",
		Method: describer.Read,
		In:     describer.Inputs{Name: "food.OrderID"},
		Out: map[describer.ResponseType]describer.ResponseObject{
			describer.ResponseOK: "accounting.Order",
		},
	}

	return describer.Service{
		Name:           "accounting",
		Rpc:            describer.GoNative,
		Dependencies:   describer.Struct,
		ConfigItems:    []describer.ConfigItem{{Name: "KafkaConsumptionRPS", Typ: "int"}},
		Endpoints:      []describer.Endpoint{getOrder},
		KafkaConsumers: []describer.TopicDesc{TopicFoodOrderUpdated},
		KafkaProducers: []describer.TopicDesc{TopicAccountingOrderPaid},
		Postgres: []describer.Postgres{{
			Name: "orderer",
			Methods: []describer.InterfaceMethod{
				{
					Name: "SaveFinishedOrder",
					Arg:  describer.InfraObject{Name: "orderID", Typ: "accounting.Order"},
					Ret:  describer.InfraObject{Name: "order", Typ: "accounting.Order"},
				},
				{
					Name: "ReadOrderByFoodID",
					Arg:  describer.InfraObject{Name: "newOrder", Typ: "food.OrderID"},
					Ret:  describer.InfraObject{Name: "order", Typ: "accounting.Order"},
				},
			},
		}},
		Clients: []describer.Client{
			{
				VarName: "pricingClient",
				IName:   "IPricingClient",
			},
		},
	}
}
