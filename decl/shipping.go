package decl

import (
	"tonky/holistic/describer"
	"tonky/holistic/typs"
)

func ShippingService() describer.ServiceV2 {
	return describer.ServiceV2{
		Tele: describer.InfraInterface{
			Interface: typs.Object3{
				Name:         "Otel",
				Module:       "tonky/holistic",
				RelativePath: []string{"infra", "tele"},
			},
			Model: typs.Object3{
				Name:         "",
				Module:       "tonky/holistic",
				RelativePath: []string{"infra", "otelinit"},
			},
		},
		/*
			Logger: describer.InfraInterface{
				Interface: typs.Object3{
					Name:         "ILogger",
					Module:       "tonky/holistic",
					RelativePath: []string{"infra", "logger"},
				},
				Model: typs.Object3{
					Name:         "Logger",
					Module:       "tonky/holistic",
					RelativePath: []string{"infra", "slogLogger"},
				},
			},
		*/
		Name:         "shipping",
		Rpc:          describer.HTTP,
		Dependencies: describer.Struct,
		ConfigItems:  []describer.ConfigItemV2{{Model: typs.Object3{Name: "KafkaConsumptionRPSLimit", Typ: typs.Int2}}},
		Endpoints: []describer.EndpointV2{
			{Name: "GetOrderByID", In: FoodOrderID, Out: FoodOrder},
		},
		KafkaConsumers: []describer.TopicDesc2{OrderPaidTopic},
		KafkaProducers: []describer.TopicDesc2{OrderShipped},
		Postgres: describer.EndpointGroups{
			{
				Name: "FoodOrderer",
				Endpoints: []describer.EndpointV2{
					{Name: "GetOrderByID", In: FoodOrderID, Out: FoodOrder},
				},
			},
		},
		/*
			Clients: []describer.InfraV2{
				{
					Name: "pricingClient",
					Model: typs.Object3{
						Kind:         typs.KindExternal,
						Name:         "PricingClient",
						Module:       "tonky/holistic",
						RelativePath: []string{"clients", "pricingClient"},
					},
				},
			},
		*/
	}
}

var FoodOrder = typs.Object3{
	Kind:         typs.KindDomain,
	Typ:          typs.Struct2,
	Name:         "Order",
	Module:       "tonky/holistic",
	RelativePath: []string{"domain", "food"},
	Fields: []typs.Object3{
		FoodOrderID,
		{Name: "Content", Typ: typs.String2, Kind: typs.KindBasic},
		{Name: "IsFinal", Typ: typs.Bool2, Kind: typs.KindBasic},
	},
}

var NewFoodOrder = typs.Object3{
	Kind:         typs.KindDomain,
	Typ:          typs.Struct2,
	Name:         "NewOrder",
	Module:       "tonky/holistic",
	RelativePath: []string{"domain", "food"},
	Fields: []typs.Object3{
		{Name: "Name", Typ: typs.String2, Kind: typs.KindBasic},
		{Name: "IsComplete", Typ: typs.Bool2, Kind: typs.KindBasic},
	},
}

var FoodOrderID = typs.Object3{
	Kind:         typs.KindDomain,
	Typ:          typs.Struct2,
	Name:         "OrderID",
	Module:       "tonky/holistic",
	RelativePath: []string{"domain", "food"},
	Fields: []typs.Object3{
		{Name: "id", Typ: typs.UUID2, Kind: typs.KindBasic},
	},
}

var ShippingOrder = typs.Object3{
	Module:       "tonky/holistic",
	Kind:         typs.KindDomain,
	Typ:          typs.Struct2,
	Name:         "Order",
	RelativePath: []string{"domain", "shipping"},
	Fields: []typs.Object3{
		FoodOrderID,
		{Name: "Address", Typ: typs.String2, Kind: typs.KindBasic},
		{Name: "Price", Typ: typs.Int2, Kind: typs.KindBasic},
	},
}

var OrderPaidTopic = describer.TopicDesc2{Name: "orderPaid", TopicName: "order.paid", Obj: FoodOrder}
var OrderShipped = describer.TopicDesc2{Name: "orderShipped", TopicName: "order.shipped", Obj: ShippingOrder}

var KafkaTopicsv2 = []describer.TopicDesc2{
	OrderPaidTopic,
	OrderShipped,
}
