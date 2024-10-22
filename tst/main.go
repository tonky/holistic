package main

import (
	"context"
	"fmt"
	"log"
	"tonky/holistic/clients"
	"tonky/holistic/domain/food"
)

func main() {
	conf := clients.Config{
		Host: "localhost",
		Port: 1234,
	}

	pc := clients.NewPizzeria(conf)

	oid, err := food.NewOrderID("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		log.Fatal("new order id err:", err)
	}

	order, err := pc.ReadOrder(context.TODO(), oid)

	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Printf("ReadOrder(%+v)\nReply: %+v\n\n", oid, order)

	newOrder := food.Order{
		ID:      oid,
		Content: "new order",
	}

	createdOrder, err := pc.CreateOrder(context.TODO(), newOrder)

	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Printf("CreateOrder(%+v)\nReply: %+v\n\n", newOrder, createdOrder)
}
