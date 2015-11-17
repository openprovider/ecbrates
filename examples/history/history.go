package main

import (
	"fmt"
	"log"

	"github.com/openprovider/ecbrates"
)

func main() {
	rates, err := ecbrates.Load()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Show history of exchange rates for EUR -> USD
	for _, r := range rates {
		if value, ok := r.Rate[ecbrates.USD].(string); ok {
			fmt.Println("Exchange rate", r.Date, ": EUR 1 -> USD", value)
		}
	}
}
