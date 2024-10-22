package app

import (
	"context"
	"fmt"
	"tonky/holistic/domain/food"
)

func NewMock() Orderer {
	return MockApp{}
}

type MockApp struct{}

func (a MockApp) ReadOrder(ctx context.Context, req food.OrderID) (food.Order, error) {
	fmt.Println("MockApp.ReadOrder", req)
	return mockReadOrder(req)
}

func (a MockApp) CreateOrder(ctx context.Context, req food.Order) (food.Order, error) {
	fmt.Println("MockApp.CrateOrder", req)
	return mockCreateOrder(req)
}
