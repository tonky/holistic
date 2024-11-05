package main

import (
	svc "tonky/holistic/services/pizzeria"
)

func main() {
	svc, err := svc.NewFromEnv()
	if err != nil {
		panic(err)
	}

	svc.Start()
}
