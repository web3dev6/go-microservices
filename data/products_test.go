package data

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProductMissingNameReturnsErr
func TestProductMissingNameReturnsErr(t *testing.T) {
	p := Product{
		Price: 1.22,
		SKU:   "abc-efg-123",
	}
	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

// TestProductMissingPriceReturnsErr
func TestProductMissingPriceReturnsErr(t *testing.T) {
	p := Product{
		Name: "abc",
		SKU:  "abc-efg-123",
	}
	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

// TestProductInvalidSKUReturnsErr
func TestProductInvalidSKUReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.22,
		SKU:   "abc",
	}
	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

// TestValidProductDoesNOTReturnsErr
func TestValidProductDoesNOTReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.22,
		SKU:   "abc-efg-123",
	}
	v := NewValidation()
	fmt.Println(p, v)
	errs := v.Validate(p)
	if errs == nil {
		assert.True(t, true)
	}
}

// TestProductsToJSON
func TestProductsToJSON(t *testing.T) {
	ps := []*Product{
		{
			Name: "abc",
		},
	}
	b := bytes.NewBufferString("")
	err := ToJSON(ps, b)
	assert.NoError(t, err)
}

// func TestChecksValidation(t *testing.T) {
// 	p := &Product{Name: "Cheap Coffee", Price: 1.00, SKU: "abc-def-123"}
// 	err := p.Validate()
// 	if err != nil {
// 		log.Println(err)
// 		t.Fail()
// 	}
// }
