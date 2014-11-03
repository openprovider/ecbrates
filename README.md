European Central Bank exchange rates
====================================

A package to get the ECB exchange rates for use with Go (golang) services

[![Build Status](https://travis-ci.org/takama/ecbrates.png?branch=master)](https://travis-ci.org/takama/ecbrates)
[![GoDoc](https://godoc.org/github.com/takama/ecbrates?status.svg)](https://godoc.org/github.com/takama/ecbrates)

### Example

```go
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
	fmt.Println("Exchange rate", r.Date, ": EUR 1 -> USD", r.Rate[ecbrates.USD])

	// Case 2: convert of 100 euros to dollars
	if value, err := r.Convert(100, ecbrates.EUR, ecbrates.USD); err == nil {
		fmt.Println("Exchange rate", r.Date, ": EUR 100.0 -> USD", value)
	}

	// Case 3: convert of 100 dollars to yens
	if value, err := r.Convert(100, ecbrates.USD, ecbrates.JPY); err == nil {
		fmt.Println("Exchange rate", r.Date, ": USD 100.0 -> JPY", value)
	}
}
```

## Author

[Igor Dolzhikov](https://github.com/takama)

## Contributors

All the contributors are welcome. If you would like to be the contributor please accept some rules.
- The pull requests will be accepted only in "develop" branch
- All modifications or additions should be tested
- Sorry, I'll not accept code with any dependency, only standard library

Thank you for your understanding!

## License

[MIT Public License](https://github.com/takama/ecbrates/blob/master/LICENSE)
