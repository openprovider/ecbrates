// Copyright 2014 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

/*
Package ecbrates 0.1.0

Example:

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
	}

European Central Bank Exchange Rates
*/
package ecbrates

import (
	"encoding/xml"
	"math"
	"net/http"
)

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
	JPY Currency = "JPY"
	BGN Currency = "BGN"
	CZK Currency = "CZK"
	DKK Currency = "DKK"
	GBP Currency = "GBP"
	HUF Currency = "HUF"
	LTL Currency = "LTL"
	PLN Currency = "PLN"
	RON Currency = "RON"
	SEK Currency = "SEK"
	CHF Currency = "CHF"
	NOK Currency = "NOK"
	HRK Currency = "HRK"
	RUB Currency = "RUB"
	TRY Currency = "TRY"
	AUD Currency = "AUD"
	BRL Currency = "BRL"
	CAD Currency = "CAD"
	CNY Currency = "CNY"
	HKD Currency = "HKD"
	IDR Currency = "IDR"
	ILS Currency = "ILS"
	INR Currency = "INR"
	KRW Currency = "KRW"
	MXN Currency = "MXN"
	MYR Currency = "MYR"
	NZD Currency = "NZD"
	PHP Currency = "PHP"
	SGD Currency = "SGD"
	THB Currency = "THB"
	ZAR Currency = "ZAR"

	rates_URL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
)

type Currency string

type Rates struct {
	Date string
	Rate map[Currency]float32
}

func New() (*Rates, error) {
	r := new(Rates)
	err := r.fetch()
	return r, err
}

// Convert value "from" -> "to"
func (r *Rates) Convert(value float32, from, to Currency) float32 {
	return round32(value*r.Rate[to]/r.Rate[from], 4)
}

// ECB XML envelope
type envelope struct {
	Data struct {
		Date  string `xml:"time,attr"`
		Rates []struct {
			Currency string  `xml:"currency,attr"`
			Rate     float32 `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

// Fetch exchange rates
func (r *Rates) fetch() error {
	r.Rate = make(map[Currency]float32)

	// a exchange rates fetched relatively the EUR currency
	r.Rate[EUR] = 1

	response, err := http.Get(rates_URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var raw envelope

	if err := xml.NewDecoder(response.Body).Decode(&raw); err != nil {
		return err
	}

	r.Date = raw.Data.Date

	for _, item := range raw.Data.Rates {
		r.Rate[Currency(item.Currency)] = item.Rate
	}

	return nil
}

func round32(x float32, prec int) float32 {
	if math.IsNaN(float64(x)) || math.IsInf(float64(x), 0) {
		return x
	}

	sign := 1.0
	if x < 0 {
		sign = -1
		x *= -1
	}

	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := float64(x) * pow
	_, frac := math.Modf(intermed)

	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return float32(rounder / pow * sign)
}
