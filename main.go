package main

import (
	"flag"
	"fmt"
	"tonky/holistic/decl"
	"tonky/holistic/generator"
)

func main() {
	genServiceName := flag.String("gen", "", "generate service code")
	flag.Parse()
	fmt.Printf("my cmd: \"%v\"\n", string(*genServiceName))

	// fmt.Printf("\n%s", pizzeria.New().Debug())

	generator.GenModels(decl.DomainModels)

	topics := decl.KafkaTopics

	generator.GenKafka(topics)

	generator.GenService(decl.PizzeriaService())
	generator.GenService(decl.AccountingService())
	generator.GenService(decl.PricingService())
	generator.GenService(decl.ShippingService())

	// pizzeria.Generate()
	// pizzeria.GenScrig()
}
