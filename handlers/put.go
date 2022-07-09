package handlers

import (
	"fmt"
	"net/http"

	"github.com/satoshi-u/go-microservices/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// UpdateProducts handles PUT requests to update products
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] Handle Products PUT ****** START ******")
	// Getting product from r.Context as middleware would have run and decoded r.Body and put product in r.Context()
	// note *** cast returned interface to data.Product
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	// invoke UpdateProduct func in package data(acts as DAL)
	product, err := data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] Product Not Found for id: ", prod.ID)
		http.Error(rw, fmt.Sprintf("Product Not Found for id: %d", prod.ID), http.StatusNotFound)
		return
	}

	// Marshal product for readable logging and log
	prodJson, _ := product.JsonMarshalProduct()
	p.l.Println("[DEBUG] Product Updated: ", string(prodJson))

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)

	p.l.Println("[DEBUG] Handle Products PUT ****** END ******")
	p.l.Println("------------------------------------------------")
}
