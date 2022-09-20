package handlers

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/satoshi-u/go-microservices/product-api/data"
)

// Products is an http.handler
type Products struct {
	l   hclog.Logger
	v   *data.Validation
	pdb *data.ProductsDB
}

// NewProducts returns a new products handler with the given logger & validator
func NewProducts(l hclog.Logger, v *data.Validation, pdb *data.ProductsDB) *Products {
	return &Products{l, v, pdb}
}

// KeyProduct to use as key when putting Product to r.Context()
type KeyProduct struct{}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("invalid path, path should be /products/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
