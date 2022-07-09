package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/satoshi-u/go-microservices/data"
)

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(rw http.ResponseWriter, r *http.Request) int {
	// Getting id from URI using gorilla mux vars
	vars := mux.Vars(r)
	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		// panic(err)
		log.Println("[ERROR] Unable to convert id from string to int")
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return -1
	}
	return id
}
