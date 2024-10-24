package main

import (
	"context"
	"fmt"
	"log"
	"tonky/holistic/clients"
	svc "tonky/holistic/services/pizzeria"
)

func main() {
	conf := clients.Config{
		Host: "localhost",
		Port: 1234,
	}

	pc := clients.NewPizzeria(conf)

	newOrder := svc.NewOrder{
		Content: "new order",
	}

	createdOrder, err := pc.CreateOrder(context.TODO(), newOrder)

	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Printf("CreateOrder(%+v)\nReply: %+v\n\n", newOrder, createdOrder)

	order, err := pc.ReadOrder(context.TODO(), createdOrder.ID)

	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Printf("ReadOrder(%+v)\nReply: %+v\n\n", createdOrder.ID, order)

}
