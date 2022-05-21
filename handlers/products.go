package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/satoshi-u/go-microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP - handler
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// log Request
	p.l.Println("Request received :::: Products Handler")
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProducts(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		p.l.Println("MethodPut")
		// expect the id in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		group := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(group) != 1 {
			p.l.Println("Invalid URI: more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(group[0]) != 2 {
			p.l.Println("Invalid URI: more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := group[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI: unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.l.Println("Got id: ", id)
		p.updateProducts(id, rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	// get products from data
	lp := data.GetProducts()
	// marshall with json.Marshal
	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	// }
	// rw.Write(d)
	// encode with json.NewEncoder
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to encode Products to json", http.StatusInternalServerError)
	}
}

// addProducts
func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")
	product := &data.Product{}
	// decode prod Json from body
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json to Product", http.StatusBadRequest)
	}
	p.l.Printf("Product: %#v", product)

	data.AddProduct(product)
}

// updateProducts
func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products")
	product := &data.Product{}
	// decode prod Json from body
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json to Product", http.StatusBadRequest)
	}
	p.l.Printf("Product: %#v", product)

	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}
}
