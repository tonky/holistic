package main

import (
	svc "tonky/holistic/services/accounting"
)

func main() {
	svc, err := svc.NewFromEnv()
	if err != nil {
		panic(err)
	}

	svc.Start()
}
