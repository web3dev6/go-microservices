package handlers

import (
	"net/http"

	"github.com/satoshi-u/go-microservices/product-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	204: noContentResponse
//  400: errorResponse
//  404: errorResponse
//  422: errorValidation

// UpdateProducts handles PUT requests to update products
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] Handle Products PUT ****** START ******")
	// As per swagger docs, header resp type : application/json
	rw.Header().Add("Content-Type", "application/json")

	// Getting product from r.Context as middleware would have run and decoded r.Body and put product in r.Context()
	// note *** cast returned interface to data.Product
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	// invoke UpdateProduct func in package data(acts as DAL)
	product, err := data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] Product Not Found for id: ", prod.ID)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// Marshal product for readable logging and log
	prodJson, _ := product.JsonMarshalProduct()
	p.l.Println("[DEBUG] Product Updated: ", string(prodJson))

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)

	p.l.Println("[DEBUG] Handle Products PUT ****** END ******")
	p.l.Println("------------------------------------------------")
	// todo : id must be required here - validation, not like AddProduct
}
