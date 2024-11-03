package main

import (
	"flag"
	"fmt"
	"tonky/holistic/decl"
	"tonky/holistic/generator"
	"tonky/holistic/generator/domain"
)

func main() {
	genServiceName := flag.String("gen", "", "generate service code")
	flag.Parse()
	fmt.Printf("my cmd: \"%v\"\n", string(*genServiceName))

	// fmt.Printf("\n%s", pizzeria.New().Debug())

	domain.Generate()

	pizzeriaDecl := decl.PizzeriaService()
	generator.GenService(pizzeriaDecl)
	generator.GenService(decl.BackofficeService())
	generator.GenService(decl.LegacyService())
	// pizzeria.Generate()
	// pizzeria.GenScrig()
}
