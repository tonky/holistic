package decl

import "tonky/holistic/typs"

var FoodOrderID = typs.Object{
	Domain: "food",
	Name:   "OrderID",
	Fields: []typs.Field{
		{Name: "ID", T: "uuid.UUID"},
	},
}

var FoodOrder = typs.Object{
	Domain: "food",
	Name:   "Order",
	Fields: []typs.Field{
		{Name: "ID", T: "food.OrderID"},
		{Name: "Content", T: "string"},
		{Name: "IsFinal", T: "bool"},
	},
}

var AccountingOrder = typs.Object{
	Domain: "accounting",
	Name:   "Order",
	Fields: []typs.Field{
		{Name: "ID", T: "food.OrderID"},
		{Name: "Cost", T: "int"},
	},
}

var ShippedOrder = typs.Object{
	Domain: "shipping",
	Name:   "Order",
	Fields: []typs.Field{
		{Name: "ID", T: "food.OrderID"},
		{Name: "ShippedAt", T: "time.Time"},
	},
}

var OrderPrice = typs.Object{
	Domain: "pricing",
	Name:   "OrderPrice",
	Fields: []typs.Field{
		{Name: "ID", T: "food.OrderID"},
		{Name: "Cost", T: "int"},
	},
}

var DomainModels = []typs.Object{FoodOrderID, FoodOrder, AccountingOrder, ShippedOrder, OrderPrice}
