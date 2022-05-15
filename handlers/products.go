package handlers

import (
	"log"
	"net/http"

	"github.com/satoshi-u/go-microservices/handlers/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// log Request
	p.l.Println("Request received @products handler")
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
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
		http.Error(rw, "Unable to encode to json", http.StatusInternalServerError)
	}
}
