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

	for currency, rate := range r.Rate {
		str, ok := rate.(string)
		if !ok {
			t.Error("Parse string error:", err)
		}
		if !currency.IsValid() {
			t.Error("Unknown currency type", currency)
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

func TestFetch90DaysExchangeRates(t *testing.T) {
	rates, err := Load()
	if err != nil {
		t.Error(err)
	}

	if len(rates) < 50 {
		t.Error("Insufficient Count of days, got", len(rates))
	}
	if len(rates) > 90 {
		t.Error("Too big Count of days, got", len(rates))
	}
	for _, item := range rates {
		if item.Date == "" {
			t.Error("Date is empty")
		}
		for currency, rate := range item.Rate {
			if !currency.IsValid() {
				t.Error("Unknown currency type", currency)
			}
			if str, ok := rate.(string); ok {
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

func TestFetchAllExchangeRates(t *testing.T) {
	rates, err := LoadAll()
	if err != nil {
		t.Error(err)
	}

	if len(rates) <= 90 {
		t.Error("Insufficient Count of days, got", len(rates))
	}
	for _, item := range rates {
		if item.Date == "" {
			t.Error("Date is empty")
		}
		for currency, rate := range item.Rate {
			if !currency.IsValid() {
				t.Error("Unknown currency type", currency)
			}
			if str, ok := rate.(string); ok {
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
