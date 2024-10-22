package app

import "tonky/holistic/domain/food"

func mockReadOrder(req food.OrderID) (food.Order, error) {
	return food.Order{ID: req, Content: "mock content"}, nil
}

func mockCreateOrder(req food.Order) (food.Order, error) {
	return req, nil
}
