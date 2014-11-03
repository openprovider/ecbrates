// Copyright 2014 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package ecbrates

import (
	"testing"
)

func TestFetchExchangeRates(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Error(err)
	}
	expected := round32(100.0*r.Rate[USD], 4)
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
