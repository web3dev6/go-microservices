package handlers

import (
	"context"
	"net/http"

	"github.com/satoshi-u/go-microservices/product-api/data"
)

// MiddlewareValidateProduct validates the product in the request and calls next if ok
func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		p.l.Debug("validation", "Extracting Product from r.Body POST|PUT")
		// As per swagger docs, header resp type : application/json
		rw.Header().Add("Content-Type", "application/json")

		product := &data.Product{}
		// Decode product from r.Body(Json)
		err := product.FromJSON(r.Body)
		if err != nil {
			p.l.Error("unable to deserialize product from r.Body", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}
		p.l.Debug("Product in r.Body: %#v", product)

		errs := p.v.Validate(product)
		if errs != nil {
			p.l.Error("error validating product in middleware", errs)
			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			// write err to rw using data.ToJSON
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}
		p.l.Debug("Product Validation Success!")

		// Put the product/productInfo in r.Context() with KeyProduct{} as key
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
