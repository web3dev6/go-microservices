package data

import (
	"log"
	"testing"
)

func TestChecksValidation(t *testing.T) {
	p := &Product{Name: "Cheap Coffee", Price: 1.00, SKU: "abc-def-123"}

	err := p.Validate()

	if err != nil {
		log.Println(err)
		t.Fail()
	}
}
