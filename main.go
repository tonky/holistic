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

	// generator.GenModels(decl.DomainModels)

	// topics := decl.KafkaTopics

	gs := generator.ServiceGen{TemplatePath: "templates", ModulePath: "tonky/holistic"}

	// gs.GenerateKafka(topics)

	// gs.Generate(decl.PizzeriaService())
	// gs.Generate(decl.AccountingService())
	// gs.Generate(decl.PricingService())
	// gs.Generate(decl.ShippingService())

	gs.GenKafka2(decl.KafkaTopicsv2)
	gs.Generate2(decl.ShippingService())
}
