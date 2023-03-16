package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/go-hclog"
	"github.com/satoshi-u/go-microservices/currency/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	Price float64 `json:"price" validate:"required,gt=0"`

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

type ProductsDB struct {
	cc          pb.CurrencyClient // not to pass by ref, since it's an interface
	log         hclog.Logger
	ratesCached map[string]float64               // cached rates
	subRClient  pb.Currency_SubscribeRatesClient // client instance for pdb
}

// New ProductsDB
func NewProductsDB(cc pb.CurrencyClient, l hclog.Logger) *ProductsDB {
	pdb := &ProductsDB{cc, l, map[string]float64{}, nil}
	go pdb.handleUpdates() // listens in background for updated rates for current client
	return pdb
}

// handleUpdates- subscribed client receives updated Rate Responses
func (pdb *ProductsDB) handleUpdates() {
	// instantiate subRClient
	subRClient, err := pdb.cc.SubscribeRates(context.Background())
	if err != nil {
		pdb.log.Error("Unable to subscribe for rates", "error", err)
		return
	}

	// save client instance in pdb
	pdb.subRClient = subRClient

	// listening in loop for rate updates,
	// if duplicate subscription request sent - handle @ gRPC Error messages in gRPC bi-directional stream - { client side }
	for {
		rr, err := subRClient.Recv() // @gRPC stream{client <- server}

		// duplicate subscription error check
		if grpcError := rr.GetError(); grpcError != nil {
			// grpcError.Code
			pdb.log.Error("error subscribing for rates", "error", err)
			continue
		}

		// valid rate-response, not any random error
		if resp := rr.GetRateResponse(); resp != nil {
			pdb.log.Info("Received updated rate from server", "dest", resp.GetDestination().String())

			if err != nil {
				pdb.log.Error("Error receiving message", "error", err)
				return
			}

			pdb.ratesCached[resp.Destination.String()] = resp.Rate
		}

	}
}

// GetProducts returns a list of products
func (pdb *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	rate, err := pdb.fetchRate(currency)
	if err != nil {
		pdb.log.Error("unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	pr := Products{}
	for _, p := range productList {
		np := *p // np is a copy, not ref
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}
	return pr, nil
}

// AddProduct adds a product to list(no err expected for now)
func (pdb *ProductsDB) AddProduct(p *Product) *Product {
	p.ID = getNextId()
	productList = append(productList, p)
	return p
}

// UpdateProduct updates an existing product in list
func (pdb *ProductsDB) UpdateProduct(p *Product) (*Product, error) {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return nil, ErrProductNotFound
	}
	// update product in db
	productList[i] = p
	return p, nil
}

// DeleteProduct deletes a product from the database
func (pdb *ProductsDB) DeleteProduct(id int) (*Product, error) {
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
func (pdb *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return productList[i], nil
	}

	rate, err := pdb.fetchRate(currency)
	if err != nil {
		pdb.log.Error("unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	np := *productList[i] // copy of product, note: product is not a deep object, flat struct
	np.Price = np.Price * rate

	return &np, nil
}

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

// helper-  get exchange rate for destination currency, base currency is EUR
func (pdb *ProductsDB) fetchRate(destination string) (float64, error) {
	// if uncommented, duplicate subscription rate request will never go thru, no gRPC bi-directional error will surface
	/*
		// if cached, return
		if r, ok := pdb.ratesCached[destination]; ok {
			return r, nil
		}
	*/

	// or get initial rate first time
	rr := &pb.RateRequest{
		Base:        pb.Currencies(pb.Currencies_value["EUR"]), // *** EUR as base always
		Destination: pb.Currencies(pb.Currencies_value[destination]),
	}
	resp, err := pdb.cc.GetRate(context.Background(), rr)
	// gRPC Error messages in Unary RPCs - at client side
	if err != nil {
		if s, ok := status.FromError(err); ok {
			// gRPC err message - yes
			metaData := s.Details()[0].(*pb.RateRequest)
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("unable to get rate from currency server { DESTINATION AND BASE CURRENCIES CANNOT BE THE SAME }, base: %s, dest: %s", metaData.Base.String(), metaData.Destination.String())
			}
			return -1, fmt.Errorf("unable to get rate from currency server, base: %s, dest: %s", metaData.Base.String(), metaData.Destination.String())
		}
		return -1, err
	}
	pdb.ratesCached[destination] = resp.Rate // update cache for first time

	// subscribe for updated rates for destination currency
	pdb.subRClient.Send(rr) // @gRPC stream{client -> server}

	return resp.Rate, err
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
