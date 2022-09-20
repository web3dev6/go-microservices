package handlers

import (
	"net/http"

	"github.com/satoshi-u/go-microservices/product-api/data"
)

// swagger:route GET /products products getProducts
//
// Returns a list of products from the database
//
//     Responses:
//       200: productsResponse
//       500: errorResponse

// GetProducts handles GET requests and returns all current products
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Handle Products GET ****** START ******")
	// As per swagger docs, header resp type : application/json
	rw.Header().Add("Content-Type", "application/json")

	// get preferred currency if it exists
	cur := r.URL.Query().Get("currency")

	// Getting products from data package
	prods, err := p.pdb.GetProducts(cur)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// Marshal products list for readable logging and log
	prodsJson, err := prods.JsonMarshalProducts()
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serialize products", "error", err)
		return
	}
	p.l.Debug("Products List: ", string(prodsJson))

	// Marshall with json.Marshal to send in ResponseWriter
	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	// 	return
	// }
	// rw.Write(d)

	// Encoding with json.NewEncoder to send in ResponseWriter
	err = prods.ToJSON(rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("unable to serialize products", "error", err)
		return
	}
	p.l.Debug("Handle Products GET ****** END ******")
	p.l.Debug("------------------------------------------------")
}

// swagger:route GET /products/{id} products getProduct
//
// Returns the product with given id from db
//
//     Responses:
//       200: productResponse
//       400: errorResponse
//       404: errorResponse
//       500: errorResponse

// GetProduct handles GET requests to return a specific product by Id
func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Handle Products GET ****** START ******")
	// As per swagger docs, header resp type : application/json
	rw.Header().Add("Content-Type", "application/json")

	// get product id from request url
	p.l.Debug("Getting product Id from url")
	id := getProductID(rw, r)
	if id == -1 {
		return
	}
	// get preferred currency if it exists
	cur := r.URL.Query().Get("currency")

	// get product from db
	p.l.Debug("Getting Product with id: ", id)
	prod, err := p.pdb.GetProductByID(id, cur)

	// handle types of errors
	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Error("product not found", "error", err)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("unable to fetch product", "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// Marshal product for readable logging and log
	prodJson, err := prod.JsonMarshalProduct()
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serialize product", "error", err)
		return
	}
	p.l.Debug("Product: ", string(prodJson))

	// write to rw using data.ToJSON
	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("unable to serialize product", "error", err)
		return
	}
	p.l.Debug("Handle Products GET ****** END ******")
	p.l.Debug("------------------------------------------------")
}
