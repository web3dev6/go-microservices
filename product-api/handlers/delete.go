package handlers

import (
	"net/http"

	"github.com/satoshi-u/go-microservices/product-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
//
// Deletes a product from the database
//
//     Responses:
//	     204: noContentResponse
//       400: errorResponse
//       404: errorResponse
//       500: errorResponse

// DeleteProducts handles DELETE requests and deletes products from the database
func (p *Products) DeleteProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Handle Products DELETE ****** START ******")
	// As per swagger docs, header resp type : application/json
	rw.Header().Add("Content-Type", "application/json")

	// get product id from request url
	p.l.Debug("Getting product Id from url")
	id := getProductID(rw, r)
	if id == -1 {
		return
	}

	p.l.Debug("Deleting in Products for id: ", id)
	// DeleteProduct IN pdb now
	product, err := p.pdb.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Debug("Product Not Found for id: ", id)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	if err != nil {
		p.l.Error("Internal server error in deleting in Products for id", id)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// Marshal product for readable logging and log
	prodJson, err := product.JsonMarshalProduct()
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serialize product", "error", err)
		return
	}
	p.l.Debug("Product Deleted: ", string(prodJson))

	// Encoding deletedProduct with json.NewEncoder to send in ResponseWriter
	// rw.Write([]byte("Product Deleted successfully"))
	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)

	p.l.Debug("Handle Products DELETE ****** END ******")
	p.l.Debug("------------------------------------------------")
}
