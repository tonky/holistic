package accountingClient

import (
	"context"
	"fmt"
	"tonky/holistic/domain/accounting"
	"tonky/holistic/domain/food"
)

func NewMock() *AccountingClientMock {
	return &AccountingClientMock{map[food.OrderID]accounting.Order{}}
}

type AccountingClientMock struct {
	Orders map[food.OrderID]accounting.Order
}

func (m *AccountingClientMock) ReadOrder(ctx context.Context, req food.OrderID) (accounting.Order, error) {
	if order, ok := m.Orders[req]; ok {
		return order, nil
	}

	return accounting.Order{}, fmt.Errorf("mock order not found")
}
