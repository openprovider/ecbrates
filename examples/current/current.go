package main

import (
	"fmt"
	"log"

	"github.com/openprovider/ecbrates"
)

func main() {
	r, err := ecbrates.New()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Case 1: get dollar rate relative to euro
	if value, ok := r.Rate[ecbrates.USD].(string); ok {
		fmt.Println("Exchange rate", r.Date, ": EUR 1 -> USD", value)
	}

	// Case 2: convert of 100 euros to dollars
	if value, err := r.Convert(100, ecbrates.EUR, ecbrates.USD); err == nil {
		fmt.Println("Exchange rate", r.Date, ": EUR 100.0 -> USD", value)
	}

	// Case 3: convert of 100 dollars to yens
	if value, err := r.Convert(100, ecbrates.USD, ecbrates.JPY); err == nil {
		fmt.Println("Exchange rate", r.Date, ": USD 100.0 -> JPY", value)
	}
}
