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
	if expected != r.Convert(100, EUR, USD) {
		t.Error("Expected rate", expected, "got", r.Convert(100, EUR, USD))
	}
}
