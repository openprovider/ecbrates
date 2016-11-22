// Copyright 2015 Openprovider Authors. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

/*
Package ecbrates 0.3.0
This package helps parse the ECB exchange rates and use it for an applications

Example 1:

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
	}

Example 2:

	package main

	import (
		"fmt"
		"log"

		"github.com/openprovider/ecbrates"
	)

	func main() {
		rates, err := ecbrates.Load() // load last 90 days
		// rates, err := ecbrates.LoadAll() // <- load ALL historical data, lots of data!
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

The European Central Bank exchange rates
*/
package ecbrates

import (
	"encoding/xml"
	"errors"
	"math"
	"net/http"
	"strconv"
)

// List of all supported currencies
const (
	AUD Currency = "AUD" // Australian Dollar (A$)
	BGN Currency = "BGN" // Bulgarian Lev (BGN)
	BRL Currency = "BRL" // Brazilian Real (R$)
	CAD Currency = "CAD" // Canadian Dollar (CA$)
	CHF Currency = "CHF" // Swiss Franc (CHF)
	CNY Currency = "CNY" // Chinese Yuan (CN¥)
	CZK Currency = "CZK" // Czech Republic Koruna (CZK)
	DKK Currency = "DKK" // Danish Krone (DKK)
	EUR Currency = "EUR" // Euro (€)
	GBP Currency = "GBP" // British Pound Sterling (£)
	HKD Currency = "HKD" // Hong Kong Dollar (HK$)
	HRK Currency = "HRK" // Croatian Kuna (HRK)
	HUF Currency = "HUF" // Hungarian Forint (HUF)
	IDR Currency = "IDR" // Indonesian Rupiah (IDR)
	ILS Currency = "ILS" // Israeli New Sheqel (₪)
	INR Currency = "INR" // Indian Rupee (Rs.)
	JPY Currency = "JPY" // Japanese Yen (¥)
	KRW Currency = "KRW" // South Korean Won (₩)
	LTL Currency = "LTL" // Lithuanian Litas (LTL)
	MXN Currency = "MXN" // Mexican Peso (MX$)
	MYR Currency = "MYR" // Malaysian Ringgit (MYR)
	NOK Currency = "NOK" // Norwegian Krone (NOK)
	NZD Currency = "NZD" // New Zealand Dollar (NZ$)
	PHP Currency = "PHP" // Philippine Peso (Php)
	PLN Currency = "PLN" // Polish Zloty (PLN)
	RON Currency = "RON" // Romanian Leu (RON)
	RUB Currency = "RUB" // Russian Ruble (RUB)
	SEK Currency = "SEK" // Swedish Krona (SEK)
	SGD Currency = "SGD" // Singapore Dollar (SGD)
	THB Currency = "THB" // Thai Baht (฿)
	TRY Currency = "TRY" // Turkish Lira (TRY)
	USD Currency = "USD" // US Dollar ($)
	ZAR Currency = "ZAR" // South African Rand (ZAR)

	// Historical currencies
	CYP Currency = "CYP"
	EEK Currency = "EEK"
	ISK Currency = "ISK"
	LVL Currency = "LVL"
	MTL Currency = "MTL"
	SIT Currency = "SIT"
	SKK Currency = "SKK"
	ROL Currency = "ROL"
	TRL Currency = "TRL"

	ratesLastURL   = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	rates90daysURL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	ratesAllURL    = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml"
)

// Currency type as a link to string
type Currency string

// Rates represent date and currency exchange rates
type Rates struct {
	Date string
	Rate map[Currency]interface{}
}

// Currencies are valid values for currency
var Currencies = []Currency{
	AUD, BGN, BRL, CAD, CHF, CNY, CZK, DKK, EUR, GBP, HKD,
	HRK, HUF, IDR, ILS, INR, JPY, KRW, LTL, MXN, MYR, NOK,
	NZD, PHP, PLN, RON, RUB, SEK, SGD, THB, TRY, USD, ZAR,

	// Historical currencies
	CYP, EEK, ISK, LVL, MTL, SIT, SKK, ROL, TRL,
}

// IsValid check Currency for valid value
func (c Currency) IsValid() bool {
	for _, value := range Currencies {
		if value == c {
			return true
		}
	}

	return false
}

// New - create a new instance of the rates and fetch a data from ECB
func New() (*Rates, error) {
	r := new(Rates)
	err := r.fetchDay()
	return r, err
}

// Load - create a new instances of the rates and fetch data for the last 90 days from ECB
func Load() ([]Rates, error) {
	return fetch90days()
}

// LoadAll - create a new instances of the rates and fetch all historical data from ECB
func LoadAll() ([]Rates, error) {
	return fetchAll()
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
	Data []struct {
		Date  string `xml:"time,attr"`
		Rates []struct {
			Currency string `xml:"currency,attr"`
			Rate     string `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

// Fetch an exchange rates
func (r *Rates) fetchDay() error {

	response, err := http.Get(ratesLastURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var raw envelope

	if err := xml.NewDecoder(response.Body).Decode(&raw); err != nil {
		return err
	}

	for _, day := range raw.Data {
		r.Rate = make(map[Currency]interface{})

		// an exchange rates fetched relatively the EUR currency
		r.Rate[EUR] = "1"

		r.Date = day.Date

		for _, item := range day.Rates {
			r.Rate[Currency(item.Currency)] = item.Rate
		}
		break
	}

	return nil
}

// Fetch a lot of exchange rates
func fetch90days() ([]Rates, error) {
	return fetchHistorical(rates90daysURL)
}

// Fetch even more exchange rates
func fetchAll() ([]Rates, error) {
	return fetchHistorical(ratesAllURL)
}

func fetchHistorical(url string) ([]Rates, error) {

	var rates []Rates

	response, err := http.Get(url)
	if err != nil {
		return rates, err
	}
	defer response.Body.Close()

	var raw envelope

	if err := xml.NewDecoder(response.Body).Decode(&raw); err != nil {
		return rates, err
	}

	for _, day := range raw.Data {

		var r Rates
		r.Rate = make(map[Currency]interface{})

		// an exchange rates fetched relatively the EUR currency
		r.Rate[EUR] = "1"

		r.Date = day.Date
		for _, item := range day.Rates {
			r.Rate[Currency(item.Currency)] = item.Rate
		}
		rates = append(rates, r)
	}
	return rates, nil
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
