package decl

import (
	"tonky/holistic/describer"
	"tonky/holistic/typs"
)

func AccountingServiceV2() describer.ServiceV2 {
	return describer.ServiceV2{
		Name:         "accountingV2",
		Rpc:          describer.HTTP,
		Dependencies: describer.Struct,
		ConfigItems:  []describer.ConfigItemV2{{Model: typs.Object3{Name: "KafkaConsumptionRPSLimit", Typ: typs.Int2}}},
		Endpoints: []describer.EndpointV2{
			{Name: "GetOrderByID", In: FoodOrderIDV2, Out: FoodOrderV2},
			{Name: "CreateOrder", In: NewFoodOrder, Out: FoodOrderV2},
		},
		KafkaConsumers: []describer.TopicDesc2{FoodOrderUpdatedTopic},
		KafkaProducers: []describer.TopicDesc2{AccountingOrderPaidTopic},
		Postgres: describer.EndpointGroups{
			{
				Name: "FoodOrderer",
				Endpoints: []describer.EndpointV2{
					{Name: "GetOrderByID", In: FoodOrderIDV2, Out: FoodOrderV2},
					{Name: "SaveNewOrder", In: NewFoodOrder, Out: FoodOrderV2},
				},
			},
		},
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
	}
}

var FoodOrderV2 = typs.Object3{
	Kind:         typs.KindDomain,
	Typ:          typs.Struct2,
	Name:         "Order",
	Module:       "tonky/holistic",
	RelativePath: []string{"domain", "foodStore"},
	Fields: []typs.Object3{
		FoodOrderIDV2,
		{Name: "Content", Typ: typs.String2, Kind: typs.KindBasic},
		{Name: "IsFinal", Typ: typs.Bool2, Kind: typs.KindBasic},
	},
}

var NewFoodOrder = typs.Object3{
	Kind:         typs.KindDomain,
	Typ:          typs.Struct2,
	Name:         "NewFoodOrder",
	Module:       "tonky/holistic",
	RelativePath: []string{"domain", "accountingV2"},
	Fields: []typs.Object3{
		{Name: "Name", Typ: typs.String2, Kind: typs.KindBasic},
		{Name: "IsComplete", Typ: typs.Bool2, Kind: typs.KindBasic},
	},
}

var FoodOrderIDV2 = typs.Object3{
	Kind:         typs.KindDomain,
	Typ:          typs.Struct2,
	Name:         "OrderID",
	Module:       "tonky/holistic",
	RelativePath: []string{"domain", "foodStore"},
	Fields: []typs.Object3{
		{Name: "id", Typ: typs.UUID2, Kind: typs.KindBasic},
	},
}

var AccountingOrderID = typs.Object3{Kind: typs.KindDomain, Typ: typs.UUID2, Name: "OrderID", RelativePath: []string{"domain", "accounting"}}
var AccountingOrderV2 = typs.Object3{Kind: typs.KindDomain, Typ: typs.Struct2, Name: "Order", RelativePath: []string{"domain", "accounting"}, Fields: []typs.Object3{AccountingOrderID}}

var FoodOrderUpdatedTopic = describer.TopicDesc2{Name: "foodOrderUpdated", TopicName: "foodStore.order.updated", Obj: FoodOrderV2}
var AccountingOrderPaidTopic = describer.TopicDesc2{Name: "accountingOrderPaid", TopicName: "accounting.order.paid", Obj: AccountingOrderV2}
