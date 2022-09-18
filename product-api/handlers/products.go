package handlers

import (
	"fmt"
	"log"

	"github.com/satoshi-u/go-microservices/currency/pb"
	"github.com/satoshi-u/go-microservices/product-api/data"
)

// Products is an http.handler
type Products struct {
	l  *log.Logger
	v  *data.Validation
	cc pb.CurrencyClient
}

// NewProducts returns a new products handler with the given logger & validator
func NewProducts(l *log.Logger, v *data.Validation, cc pb.CurrencyClient) *Products {
	return &Products{l, v, cc}
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

// ServeHTTP - handler **********************************************************************
// func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
// 	// log Request
// 	p.l.Println("Request received :::: Products Handler")
// 	if r.Method == http.MethodGet {
// 		p.getProducts(rw, r)
// 		return
// 	}
// 	if r.Method == http.MethodPost {
// 		p.addProducts(rw, r)
// 		return
// 	}
// 	if r.Method == http.MethodPut {
// 		p.l.Println("MethodPut")
// 		// expect the id in the URI
// 		reg := regexp.MustCompile(`/([0-9]+)`)
// 		group := reg.FindAllStringSubmatch(r.URL.Path, -1)
// 		if len(group) != 1 {
// 			p.l.Println("Invalid URI: more than one id")
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}
// 		if len(group[0]) != 2 {
// 			p.l.Println("Invalid URI: more than one capture group")
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}
// 		idString := group[0][1]
// 		id, err := strconv.Atoi(idString)
// 		if err != nil {
// 			p.l.Println("Invalid URI: unable to convert to number", idString)
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		p.l.Println("Got id: ", id)
// 		p.updateProducts(id, rw, r)
// 		return
// 	}

// 	// catch all
// 	rw.WriteHeader(http.StatusMethodNotAllowed)
// }
