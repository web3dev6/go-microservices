package handlers

import (
	"fmt"
	"net/http"

	"github.com/satoshi-u/go-microservices/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
//
// Deletes a product from the database
//
//     Responses:
//       201: noContentResponse
//       404: errorResponse
//       501: errorResponse

// DeleteProducts handles DELETE requests and deletes products from the database
func (p *Products) DeleteProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] Handle Products DELETE ****** START ******")
	// get product id from request url
	p.l.Println("[DEBUG] Getting product Id from url")
	id := getProductID(rw, r)
	p.l.Println("[DEBUG] Deleting in Products for id: ", id)

	// DeleteProduct func in package data(acts as DAL)
	product, err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] Product Not Found for id: ", id)
		http.Error(rw, fmt.Sprintf("Product Not Found for id: %d", id), http.StatusNotFound)
		return
	}
	if err != nil {
		p.l.Println("[ERROR] Internal server error in deleting in Products for id", id)
		http.Error(rw, fmt.Sprintf("Internal server error in deleting in Products for id: %d", id), http.StatusInternalServerError)
		return
	}

	// Marshal product for readable logging and log
	prodJson, _ := product.JsonMarshalProduct()
	p.l.Println("[DEBUG] Product Deleted: ", string(prodJson))

	// Encoding deletedProduct with json.NewEncoder to send in ResponseWriter
	// rw.Write([]byte("Product Deleted successfully"))
	err = product.ToJSON(rw)
	if err != nil {
		p.l.Println("[ERROR] Unable to encode Product to json")
		http.Error(rw, "Unable to encode Product to json", http.StatusInternalServerError)
		return
	}
	p.l.Println("[DEBUG] Handle Products DELETE ****** END ******")
	p.l.Println("------------------------------------------------")
}
