package main

import (
	"log"
	"testing"

	"github.com/satoshi-u/go-microservices/product-api/sdk/client"
	"github.com/satoshi-u/go-microservices/product-api/sdk/client/products"
	"github.com/satoshi-u/go-microservices/product-api/sdk/models"
)

// To clean test-cache, run : go clean -testcache
// todo: all err cases, tested with curl, need to add test-funcs here as well

func TestClientForGetProducts(t *testing.T) {
	// c := client.Default
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)
	params := products.NewGetProductsParams()
	prods, err := c.Products.GetProducts(params)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%#v", prods.GetPayload()[0])
	log.Printf("%#v", prods.GetPayload()[1])
	// for getting logs
	// t.Fail()
}

func TestClientForGetProduct(t *testing.T) {
	// c := client.Default
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)
	params := products.NewGetProductParams()
	params.ID = 1
	prod, err := c.Products.GetProduct(params)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%#v", prod.GetPayload())
	// for getting logs
	// t.Fail()
}

func TestClientForAddProducts(t *testing.T) {
	// c := client.Default
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)
	params := products.NewCreateProductParams()
	prodName := "mango-shake"
	prodDesc := "mango & milk"
	prodPrice := float32(6.50)
	prodSKU := "prod-bev-000"
	params.WithDefaults().SetBody(&models.Product{Name: &prodName, Description: prodDesc, Price: &prodPrice, SKU: &prodSKU})
	// todo : ensure the product is already not there in CreateProduct (duplicates)
	prod, err := c.Products.CreateProduct(params)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%#v", prod.GetPayload())
	// for getting logs
	// t.Fail()
}

func TestClientForUpdateProducts(t *testing.T) {
	// c := client.Default
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)
	params := products.NewUpdateProductParams()
	prodId := 3
	prodName := "mango-banana-shake"
	prodDesc := "mango & banana & milk mix"
	prodPrice := float32(7.50)
	prodSKU := "prod-bev-003"
	params.WithDefaults().SetBody(&models.Product{ID: int64(prodId), Name: &prodName, Description: prodDesc, Price: &prodPrice, SKU: &prodSKU})
	prodUpdated, err := c.Products.UpdateProduct(params)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("resp %#v", prodUpdated.Error())
	// for getting logs
	// t.Fail()
}

func TestClientForDeleteProducts(t *testing.T) {
	// c := client.Default
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)
	params := products.NewDeleteProductParams()
	params.ID = 3
	prodDeleted, err := c.Products.DeleteProduct(params)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("resp %#v", prodDeleted.Error())
	// for getting logs
	// t.Fail()
}
