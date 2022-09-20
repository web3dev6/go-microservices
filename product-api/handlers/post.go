package handlers

import (
	"net/http"

	"github.com/satoshi-u/go-microservices/product-api/data"
)

// swagger:route POST /products products createProduct
// Creates a new product
//
// responses:
//	200: productResponse
//  400: errorResponse
//  422: errorValidation
//  500: errorResponse

// AddProducts handles POST requests to add new products
func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Handle Products POST ****** START ******")
	// As per swagger docs, header resp type : application/json
	rw.Header().Add("Content-Type", "application/json")

	// Getting product from r.Context as middleware would have run and decoded r.Body and put product in r.Context()
	// note *** cast returned interface to data.Product
	product := r.Context().Value(KeyProduct{}).(*data.Product)

	// invoke AddProduct func in package data(acts as DAL)
	product = p.pdb.AddProduct(product)

	// Marshal product for readable logging and log
	prodJson, err := product.JsonMarshalProduct()
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serialize product", "error", err)
		return
	}
	p.l.Debug("Product Added: ", string(prodJson))

	// Encoding with json.NewEncoder to send in ResponseWriter
	// rw.Write([]byte("Product Added successfully"))

	// encode product to json
	err = product.ToJSON(rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("unable to serialize product", "error", err)
		return
	}
	p.l.Debug("Handle Products POST ****** END ******")
	p.l.Debug("------------------------------------------------")
}
