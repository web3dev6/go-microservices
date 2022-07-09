package handlers

import (
	"context"
	"net/http"

	"github.com/satoshi-u/go-microservices/data"
)

// MiddlewareValidateProduct validates the product in the request and calls next if ok
func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		p.l.Println("[DEBUG] MiddlewareValidateProduct:- *Extracting Product from r.Body POST|PUT")
		product := &data.Product{}

		// Decode product from r.Body(Json)
		err := product.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product from r.Body in middleware", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}
		p.l.Printf("[DEBUG] Product from r.Body: %#v", product)

		// validate the product
		// err = product.Validate()
		// if err != nil {
		// 	p.l.Println("[ERROR] validating product in middleware", err)
		// 	http.Error(rw, fmt.Sprintf("Error validating Product: %s", err), http.StatusBadRequest)
		// 	return
		// }
		errs := p.v.Validate(product)
		if errs != nil {
			p.l.Println("[ERROR] validating product in middleware", err)
			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			// write err to rw using data.ToJSON
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}
		p.l.Println("[DEBUG] Product Validation:- *Success")

		// Put the product/productInfo in r.Context() with KeyProduct{} as key
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
