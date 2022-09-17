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
	p.l.Println("[DEBUG] Handle Products GET ****** START ******")
	// As per swagger docs, header resp type : application/json
	rw.Header().Add("Content-Type", "application/json")

	// Getting products from data package
	prods := data.GetProducts()

	// Marshal products list for readable logging and log
	prodsJson, _ := prods.JsonMarshalProducts()
	p.l.Println("[DEBUG] Products List: ", string(prodsJson))

	// Marshall with json.Marshal to send in ResponseWriter
	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	// 	return
	// }
	// rw.Write(d)

	// Encoding with json.NewEncoder to send in ResponseWriter
	err := prods.ToJSON(rw)
	if err != nil {
		p.l.Println("[ERROR] Unable to encode Products to json")
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	p.l.Println("[DEBUG] Handle Products GET ****** END ******")
	p.l.Println("------------------------------------------------")
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
	p.l.Println("[DEBUG] Handle Products GET ****** START ******")
	// As per swagger docs, header resp type : application/json
	rw.Header().Add("Content-Type", "application/json")

	// get product id from request url
	p.l.Println("[DEBUG] Getting product Id from url")
	id := getProductID(rw, r)
	if id == -1 {
		return
	}

	// get product from db
	p.l.Println("[DEBUG] Getting Product with id: ", id)
	prod, err := data.GetProductByID(id)

	// handle types of errors
	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// Marshal product for readable logging and log
	prodJson, _ := prod.JsonMarshalProduct()
	p.l.Println("[DEBUG] Product: ", string(prodJson))

	// write to rw using data.ToJSON
	err = data.ToJSON(prod, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product to response", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	p.l.Println("[DEBUG] Handle Products GET ****** END ******")
	p.l.Println("------------------------------------------------")
}
