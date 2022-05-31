package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// FromJSON : when adding/updating a product, used in MiddlewareValidateProduct
func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	err := d.Decode(p)
	if err != nil {
		log.Printf("Unable to decode from Json for product with id{%d}, err: %v", p.ID, err)
		return err
	}
	return nil
}

// ToJSON : when adding/updating a product, used to write back as http response
func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	err := e.Encode(p)
	if err != nil {
		log.Printf("Unable to encode to Json for product with id{%d}, err: %v", p.ID, err)
		return err
	}
	return nil
}

// JsonMarshalProduct: marshals product into json, used for logging purposes
func (p *Product) JsonMarshalProduct() ([]byte, error) {
	product, err := json.Marshal(p)
	if err != nil {
		log.Printf("Unable to marshal to json for product with id{%d}, err: %v", p.ID, err)
		return nil, err
	}
	return product, nil
}

// Validate : when adding/updating a product, used in MiddlewareValidateProduct
func (p *Product) Validate() error {
	log.Println("Validate:- *Validating Product from r.Body **POST|PUT")
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	err := validate.Struct(p)
	if err != nil {
		log.Printf("error in validation for product with id{%d}, err: %v", p.ID, err)
		// validationErrors := err.(validator.ValidationErrors)
		// log.Println("validationErrors: ", validationErrors)
		return err
	}
	return nil
	// return validate.Struct(p)
}

// validateSKU : custom validation for sku field with regex
func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-asdf-1234
	regex := regexp.MustCompile(`[a-z]+-[a-z]+-[0-9]+`)
	matches := regex.FindAllString(fl.Field().String(), -1)
	return len(matches) == 1
}

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
		log.Printf("Unable to encode to Json for Products, err: %v", err)
		return err
	}
	return nil
}

// JsonMarshalProducts: marshals products list into json, used for logging purposes
func (p *Products) JsonMarshalProducts() ([]byte, error) {
	products, err := json.Marshal(p)
	if err != nil {
		log.Printf("Unable to marshal to json for products, err: %v", err)
		return nil, err
	}
	return products, nil
}

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

// AddProduct adds a product to list(no err expected for now)
func AddProduct(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}

// UpdateProduct updates an existing product in list
func UpdateProduct(id int, pInfo *Product) (*Product, error) {
	fp, pos, err := findProduct(id)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("id: %d\n", id)
	// fmt.Printf("Product: %#v\n", p)
	if fp.ID == id {
		pInfo.ID = id
		productList[pos] = pInfo
	}
	return pInfo, nil
}

// getNextId calculates ID for a new product to be added
func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var ErrProductNotFound = fmt.Errorf("Product not found")

// getNextId calculates ID for a new product to be added
func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
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
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "prod-bev-002",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
