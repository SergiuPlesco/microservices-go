package data

import (
	"testing"
)

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "V60",
		Price: 5.00,
		SKU:   "asdf-fd-asdf",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
