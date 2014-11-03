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
	str, ok := r.Rate[USD].(string)
	if !ok {
		t.Error("Parse string error:", err)
	}
	v, err := strconv.ParseFloat(str, 32)
	if !ok {
		t.Error("Parse float error:", err)
	}
	expected := round64(100.0*round64(v, 4), 4)
	value, err := r.Convert(100, EUR, USD)
	if err != nil {
		t.Error("Converting error:", err)
	}
	if expected != value {
		t.Error("Expected rate", expected, "got", value)
	}
	if _, err = r.Convert(100, Currency("XXX"), USD); err == nil {
		t.Error("Expected error, got nil")
	}
	if _, err = r.Convert(100, USD, Currency("XXX")); err == nil {
		t.Error("Expected error, got nil")
	}
	if _, err = r.Convert(100, Currency("XXX"), Currency("XXX")); err == nil {
		t.Error("Expected error, got nil")
	}
}
