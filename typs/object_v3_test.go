package typs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var domainAssortmentMealID = Object3{
	Typ:          UUID2,
	Kind:         KindDomain,
	Name:         "ID",
	ImportAlias:  "MealID",
	RelativePath: []string{"domain", "assortment"},
	Module:       "tonky/holistic",
}

var domainAssortmentMeal = Object3{
	Typ:          Struct2,
	Kind:         KindDomain,
	Name:         "Meal",
	RelativePath: []string{"domain", "assortment"},
	Module:       "tonky/holistic",
	Fields: []Object3{
		domainAssortmentMealID,
		{
			Name: "PublishedAt",
			Typ:  Time2,
			Kind: KindBasic,
		},
	},
}

var domainFoodOrderID = Object3{
	Typ:          UUID2,
	Kind:         KindDomain,
	Name:         "ID",
	ImportAlias:  "OrderID",
	RelativePath: []string{"domain", "order"},
	Module:       "tonky/holistic",
}

var domainFoodOrder = Object3{
	Typ:          Struct2,
	Kind:         KindDomain,
	Name:         "Order",
	RelativePath: []string{"domain", "order"},
	Module:       "tonky/holistic",
	Fields:       []Object3{domainAssortmentMealID, domainFoodOrderID},
}

var domainOrderDelivery = Object3{
	Typ:          Struct2,
	Kind:         KindDomain,
	Name:         "Delivery",
	RelativePath: []string{"domain", "delivery"},
	Module:       "tonky/holistic",
	Fields: []Object3{
		{
			Name: "DeliveryID",
			Typ:  UUID2,
			Kind: KindBasic,
		},
		domainAssortmentMealID,
		domainFoodOrderID,
		{
			Name: "PackagingDuration",
			Typ:  Duration2,
			Kind: KindBasic,
		},
		{
			Name: "DispatchedAt",
			Typ:  Time2,
			Kind: KindBasic,
		},
		{
			Name: "DeliveeredAt",
			Typ:  Time2,
			Kind: KindBasic,
		},
		{
			Name: "DeliveryAttempts",
			Typ:  Int2,
			Kind: KindBuiltIn,
		},
	},
}

func TestImports(t *testing.T) {
	none := []string{}

	require.Equal(t, []string{"github.com/google/uuid"}, domainFoodOrderID.AbsImportsAll(domainFoodOrderID))

	require.Equal(t, none, domainFoodOrderID.AbsImportsAll(domainFoodOrder))
	require.Equal(t, none, domainFoodOrder.AbsImportsAll(domainFoodOrderID))
	require.Equal(t, []string{"tonky/holistic/domain/order"}, domainFoodOrder.AbsImportsAll(domainAssortmentMeal))
	require.Equal(t, []string{"tonky/holistic/domain/assortment"}, domainAssortmentMeal.AbsImportsAll(domainFoodOrder))

	require.Equal(t, []string{"tonky/holistic/domain/delivery"}, domainOrderDelivery.AbsImportsAll(domainAssortmentMeal))

	modelSelfImports := []string{"github.com/google/uuid", "tonky/holistic/domain/assortment", "tonky/holistic/domain/order", "time"}
	require.Equal(t, modelSelfImports, domainOrderDelivery.AbsImports(domainOrderDelivery))
}

// func TestStructGen(t *testing.T) { }
