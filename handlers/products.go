package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/satoshi-u/go-microservices/data"
)

type Products struct {
	l *log.Logger
}

// NewProduct : new handler init
func NewProduct(l *log.Logger) *Products {
	return &Products{l}
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

// GetProducts **********************************************************************
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Products GET ****** START ******")
	// Getting products from data package
	lp := data.GetProducts()
	// Marshal products list for readable logging and log
	lpjson, _ := lp.JsonMarshalProducts()
	p.l.Println("Products List: ", string(lpjson))
	// Marshall with json.Marshal to send in ResponseWriter
	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	// 	return
	// }
	// rw.Write(d)
	// Encoding with json.NewEncoder to send in ResponseWriter
	err := lp.ToJSON(rw)
	if err != nil {
		p.l.Println("[ERROR] Unable to encode Products to json")
		http.Error(rw, "Unable to encode Products to json", http.StatusInternalServerError)
		return
	}
	p.l.Println("Handle Products GET ****** END ******")
}

// AddProducts **********************************************************************
func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Products POST ****** START ******")
	// Getting product from r.Context as middleware would have run and decoded r.Body and put product in r.Context()
	// note *** cast returned interface to data.Product
	product := r.Context().Value(KeyProduct{}).(*data.Product)
	// invoke AddProduct func in package data(acts as DAL)
	data.AddProduct(product)
	// Marshal product for readable logging and log
	pjson, _ := product.JsonMarshalProduct()
	p.l.Println("Product Added: ", string(pjson))
	// Encoding with json.NewEncoder to send in ResponseWriter
	// rw.Write([]byte("Product Added successfully"))
	err := product.ToJSON(rw)
	if err != nil {
		p.l.Println("[ERROR] Unable to encode Product to json")
		http.Error(rw, "Unable to encode Product to json", http.StatusInternalServerError)
		return
	}
	p.l.Println("Handle Products POST ****** END ******")
}

// UpdateProducts **********************************************************************
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Products PUT ****** START ******")
	// Getting id from URI using gorilla mux vars
	vars := mux.Vars(r)
	// p.l.Println("mux.Vars PUT Products", vars)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.l.Println("[ERROR] Unable to convert id from string to int")
		http.Error(rw, "Unable to convert id from string to int", http.StatusBadRequest)
		return
	}
	p.l.Println("Updating Products for id: ", id)
	// Getting product from r.Context as middleware would have run and decoded r.Body and put product in r.Context()
	// note *** cast returned interface to data.Product
	productInfo := r.Context().Value(KeyProduct{}).(*data.Product)
	// UpdateProduct func in package data(acts as DAL)
	product, err := data.UpdateProduct(id, productInfo)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] Product Not Found for id: ", id)
		http.Error(rw, fmt.Sprintf("Product Not Found for id: %d", id), http.StatusNotFound)
		return
	}
	if err != nil {
		p.l.Println("[ERROR] Internal server error in updating Products for id", id)
		http.Error(rw, fmt.Sprintf("Internal server error in updating Products for id: %d", id), http.StatusInternalServerError)
		return
	}
	// Marshal product for readable logging and log
	pjson, _ := product.JsonMarshalProduct()
	p.l.Println("Product Updated: ", string(pjson))
	// Encoding updatedProduct with json.NewEncoder to send in ResponseWriter
	// rw.Write([]byte("Product Updated successfully"))
	err = product.ToJSON(rw)
	if err != nil {
		p.l.Println("[ERROR] Unable to encode Product to json")
		http.Error(rw, "Unable to encode Product to json", http.StatusInternalServerError)
		return
	}
	p.l.Println("Handle Products PUT ****** END ******")
}

// KeyProduct to use as key when putting Product to r.Context()
type KeyProduct struct{}

// MiddlewareValidateProduct **********************************************************************
// validates/extracts Product from r.Body(Json) and puts in r.Context before handler code runs for a route
func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		p.l.Println("MiddlewareValidateProduct:- *Extracting Product Info from r.Body POST|PUT")
		product := &data.Product{}
		// Decode productInfo from r.Body(Json)
		err := product.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product from r.Body in middleware", err)
			http.Error(rw, "Unable to unmarshal json to Product", http.StatusBadRequest)
			return
		}
		p.l.Printf("Product Info from r.Body: %#v", product)

		// validate the product
		err = product.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product in middleware", err)
			http.Error(rw, fmt.Sprintf("Error validating Product: %s", err), http.StatusBadRequest)
			return
		}
		p.l.Println("Product Validation:- *Success")

		// Put the product/productInfo in r.Context() with KeyProduct{} as key
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
