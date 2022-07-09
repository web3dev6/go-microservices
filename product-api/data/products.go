package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}

// FromJSON : when adding/updating a product, used in MiddlewareValidateProduct
func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	err := d.Decode(p)
	if err != nil {
		log.Printf("[ERROR] Unable to decode from Json for product with id{%d}, err: %v", p.ID, err)
		return err
	}
	return nil
}

// ToJSON : when adding/updating a product, used to write back as http response
func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	err := e.Encode(p)
	if err != nil {
		log.Printf("[ERROR] Unable to encode to Json for product with id{%d}, err: %v", p.ID, err)
		return err
	}
	return nil
}

// JsonMarshalProduct: marshals product into json, used for logging purposes
func (p *Product) JsonMarshalProduct() ([]byte, error) {
	product, err := json.Marshal(p)
	if err != nil {
		log.Printf("[ERROR] Unable to marshal to json for product with id{%d}, err: %v", p.ID, err)
		return nil, err
	}
	return product, nil
}

// Validate : when adding/updating a product, used in MiddlewareValidateProduct
// func (p *Product) Validate() error {
// 	log.Println("Validate:- *Validating Product from r.Body **POST|PUT")
// 	validate := validator.New()
// 	validate.RegisterValidation("sku", validateSKU)
// 	err := validate.Struct(p)
// 	if err != nil {
// 		log.Printf("[ERROR] error in validation for product with id{%d}, err: %v", p.ID, err)
// 		// validationErrors := err.(validator.ValidationErrors)
// 		// log.Println("validationErrors: ", validationErrors)
// 		return err
// 	}
// 	return nil
// 	// return validate.Struct(p)
// }

// validateSKU : custom validation for sku field with regex
// func validateSKU(fl validator.FieldLevel) bool {
// 	// sku is of format abc-asdf-1234
// 	regex := regexp.MustCompile(`[a-z]+-[a-z]+-[0-9]+`)
// 	matches := regex.FindAllString(fl.Field().String(), -1)
// 	return len(matches) == 1
// }

// Products is a collection of Product
type Products []*Product

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
// Products ToJSON : used in getting all products GET
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	err := e.Encode(p)
	if err != nil {
		log.Printf("[ERROR] Unable to encode to Json for Products, err: %v", err)
		return err
	}
	return nil
}

// JsonMarshalProducts: marshals products list into json, used for logging purposes
func (p *Products) JsonMarshalProducts() ([]byte, error) {
	products, err := json.Marshal(p)
	if err != nil {
		log.Printf("[ERROR] Unable to marshal to json for products, err: %v", err)
		return nil, err
	}
	return products, nil
}

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

// AddProduct adds a product to list(no err expected for now)
func AddProduct(p *Product) *Product {
	p.ID = getNextId()
	productList = append(productList, p)
	return p
}

// UpdateProduct updates an existing product in list
func UpdateProduct(p *Product) (*Product, error) {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return nil, ErrProductNotFound
	}
	// update product in db
	productList[i] = p
	return p, nil
}

// DeleteProduct deletes a product from the database
func DeleteProduct(id int) (*Product, error) {
	i := findIndexByProductID(id)
	// log.Printf("[DEBUG] ************************** i : %v", i)
	if i == -1 {
		return nil, ErrProductNotFound
	}
	pdel := productList[i]
	// Remove the product at index i from productList.
	copy(productList[i:], productList[i+1:])       // Shift productList[i+1:] left one index.
	productList[len(productList)-1] = nil          // Erase last element (write zero value).
	productList = productList[:len(productList)-1] // Truncate slice.
	// return deleted product
	return pdel, nil
}

// getNextId calculates ID for a new product to be added
func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var ErrProductNotFound = fmt.Errorf("Product not found")

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}
	return productList[i], nil
}

// findProduct finds the product with given id
// func findProduct(id int) (*Product, int, error) {
// 	for i, p := range productList {
// 		if p.ID == id {
// 			return p, i, nil
// 		}
// 	}
// 	return nil, -1, ErrProductNotFound
// }

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

// productList is a hard coded list of products for this
// example data source
var productList = Products{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "prod-bev-001",
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "prod-bev-002",
	},
}
