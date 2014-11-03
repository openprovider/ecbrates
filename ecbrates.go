// Copyright 2014 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

/*
Package ecbrates 0.1.4

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

		// Case 1: get dollar rate relative to euro
		if value, ok := r.Rate[ecbrates.USD].(string); ok {
			fmt.Println("Exchange rate", r.Date, ": EUR 1 -> USD", value)
		}

		// Case 2: convert of 100 euros to dollars
		if value, err := r.Convert(100, ecbrates.EUR, ecbrates.USD); err == nil {
			fmt.Println("Exchange rate", r.Date, ": EUR 100.0 -> USD", value)
		}
	}

European Central Bank exchange rates
*/
package ecbrates

import (
	"encoding/xml"
	"errors"
	"math"
	"net/http"
	"strconv"
)

// Links to all supported currencies
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

	ratesURL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
)

// Currency type as a link to string
type Currency string

// Rates represent date and currency exchange rates
type Rates struct {
	Date string
	Rate map[Currency]interface{}
}

// New - create a new instance of the rates and fetch a data from ECB
func New() (*Rates, error) {
	r := new(Rates)
	err := r.fetch()
	return r, err
}

// Convert a value "from" one Currency -> "to" other Currency
func (r *Rates) Convert(value float64, from, to Currency) (float64, error) {
	if r.Rate[to] == nil || r.Rate[from] == nil {
		return 0, errors.New("Perhaps one of the values ​​of currencies does not exist")
	}
	errorMessage := "Perhaps one of the values ​​of currencies could not parsed correctly"
	strFrom, okFrom := r.Rate[from].(string)
	strTo, okTo := r.Rate[to].(string)
	if !okFrom || !okTo {
		return 0, errors.New(errorMessage)
	}
	vFrom, err := strconv.ParseFloat(strFrom, 32)
	if err != nil {
		return 0, errors.New(errorMessage)
	}
	vTo, err := strconv.ParseFloat(strTo, 32)
	if err != nil {
		return 0, errors.New(errorMessage)
	}
	return round64(value*round64(vTo, 4)/round64(vFrom, 4), 4), nil

}

// ECB XML envelope
type envelope struct {
	Data struct {
		Date  string `xml:"time,attr"`
		Rates []struct {
			Currency string `xml:"currency,attr"`
			Rate     string `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

// Fetch an exchange rates
func (r *Rates) fetch() error {
	r.Rate = make(map[Currency]interface{})

	// an exchange rates fetched relatively the EUR currency
	r.Rate[EUR] = "1"

	response, err := http.Get(ratesURL)
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

func round64(x float64, prec int) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}

	sign := 1.0
	if x < 0 {
		sign = -1
		x *= -1
	}

	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)

	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow * sign
}
