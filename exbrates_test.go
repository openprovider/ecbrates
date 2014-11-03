// Copyright 2014 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package ecbrates

import (
	"strconv"
	"testing"
)

func TestFetchExchangeRates(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Error(err)
	}
	if r.Date == "" {
		t.Error("Date is empty")
	}
	if len(r.Rate) != len(Currencies) {
		t.Error("Insufficient count of rates, got", len(r.Rate), "for", r.Date)
	}

	for _, currency := range Currencies {
		str, ok := r.Rate[currency].(string)
		if !ok {
			t.Error("Parse string error:", err)
		}
		v, err := strconv.ParseFloat(str, 32)
		if !ok {
			t.Error("Parse float error:", err)
		}
		expected := round64(100.0*round64(v, 4), 4)
		value, err := r.Convert(100, EUR, currency)
		if err != nil {
			t.Error("Converting error:", err)
		}
		if expected != value {
			t.Error("Expected rate", expected, "got", value)
		}
	}
	if _, err = r.Convert(100, Currency("XXX"), EUR); err == nil {
		t.Error("Expected error, got nil")
	}
	if _, err = r.Convert(100, EUR, Currency("XXX")); err == nil {
		t.Error("Expected error, got nil")
	}
	if _, err = r.Convert(100, Currency("XXX"), Currency("XXX")); err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestFetchAllExchangeRates(t *testing.T) {
	rates, err := Load()
	if err != nil {
		t.Error(err)
	}

	if len(rates) < 50 {
		t.Error("Insufficient Count of days, got", len(rates))
	}
	for _, item := range rates {
		if item.Date == "" {
			t.Error("Date is empty")
		}
		if len(item.Rate) != len(Currencies) {
			t.Error("Day:", item.Date, "Insufficient count of rates, got", len(item.Rate))
		}
		for _, currency := range Currencies {
			if str, ok := item.Rate[currency].(string); ok {
				if v, err := strconv.ParseFloat(str, 32); err == nil {
					if v == 0 {
						t.Error("Day:", item.Date, "Zero rate for", currency)
					}
				} else {
					t.Error(err)
				}
			} else {
				t.Error("Parse rate to string unsuccessful")
			}
		}
	}
}
