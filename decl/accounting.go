package decl

import "tonky/holistic/generator/services"

func AccountingService() services.Service {
	getOrder := services.Endpoint{
		Name:   "order",
		Method: services.Read,
		In:     services.Inputs{Name: "food.OrderID"},
		Out: map[services.ResponseType]services.ResponseObject{
			services.ResponseOK: "accounting.Order",
		},
	}

	return services.Service{
		Name:           "accounting",
		Rpc:            services.GoNative,
		ConfigItems:    []services.ConfigItem{{Name: "KafkaConsumptionRPS", Typ: "int"}},
		Endpoints:      []services.Endpoint{getOrder},
		KafkaConsumers: []services.TopicDesc{services.TopicFoodOrderUpdated},
		KafkaProducers: []services.TopicDesc{services.TopicAccountingOrderPaid},
		Interfaces: []services.JustInterface{{
			Name:   "AccountOrdersRepoReader",
			Struct: "OrdersRepository",
			Deps: map[string]services.FQImport{
				"logger": {
					RelPath: "infra/logger",
					Model:   "Slog",
				},
			},
			Methods: []services.InterfaceMethod{
				{
					Name: "SaveOrder",
					Arg:  services.InfraObject{Name: "order", Typ: "NewOrder"},
					Ret:  services.InfraObject{Name: "order", Typ: "accounting.Order"},
				},
				{
					Name: "ReadOrderByFoodID",
					Arg:  services.InfraObject{Name: "orderID", Typ: "food.OrderID"},
					Ret:  services.InfraObject{Name: "order", Typ: "accounting.Order"},
				},
			},
		}},
		Clients: []services.Client{
			{
				VarName: "pricingClient",
				IName:   "IPricingClient",
			},
		},
	}
}
