package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products is a collection of Product
type Products []*Product

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
// Products ToJSON : when getting all products
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// product FromJSON : when adding a product
func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

// AddProduct
func AddProduct(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}

// UpdateProduct
func UpdateProduct(id int, p *Product) error {
	fp, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	// fmt.Printf("id: %d\n", id)
	// fmt.Printf("Product: %#v\n", p)
	if fp.ID == id {
		p.ID = id
		productList[pos] = p
	}
	return nil
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
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
