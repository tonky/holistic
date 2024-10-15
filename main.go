package main

import (
	"flag"
	"fmt"
	"tonky/holistic/domain"
	"tonky/holistic/services/pizzeria"
)

func main() {
	genServiceName := flag.String("gen", "", "generate service code")
	flag.Parse()
	fmt.Printf("my cmd: \"%v\"\n", string(*genServiceName))

	fmt.Printf("\n%s", pizzeria.New().Debug())

	domain.Generate()
	// pizzeria.Generate()
	pizzeria.GenScrig()
}
