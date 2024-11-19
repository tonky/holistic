package main

import "tonky/holistic/services/accountingV2"

func main() {
	s, err := accountingV2.NewFromEnv()
	if err != nil {
		panic(err)
	}

	s.Start()
}
