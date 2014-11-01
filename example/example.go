package main

import (
	"fmt"
	"log"

	"github.com/takama/ecbrates"
)

func main() {
	r, err := ecbrates.New()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Variant 1: get dollar rate relative to euro
	fmt.Println("Exchange rate", r.Date, ": EUR 1 -> USD", r.Rate[ecbrates.USD])

	// Variant 2: convert of 100 euros to dollars
	fmt.Println(
		"Exchange rate", r.Date,
		": EUR 100.0 -> USD", r.Convert(100, ecbrates.EUR, ecbrates.USD),
	)

	// Variant 3: convert of 100 dollars to yens
	fmt.Println(
		"Exchange rate", r.Date,
		": USD 100.0 -> JPY", r.Convert(100, ecbrates.USD, ecbrates.JPY),
	)
}
