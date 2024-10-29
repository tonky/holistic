package accounting

import (
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"
)

type AccountingOrdersRepositoryMemory struct {
	ordersByFoodID map[food.OrderID]accounting.Order
}

func NewAccountingOrdersRepositoryMemory() *AccountingOrdersRepositoryMemory {
	return &AccountingOrdersRepositoryMemory{
		ordersByFoodID: make(map[food.OrderID]accounting.Order),
	}
}
