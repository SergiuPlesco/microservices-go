package main

import (
	"fmt"
	"testing"

	"github.com/SergiuPlesco/microservices-go/gorilla/sdk/client"
	"github.com/SergiuPlesco/microservices-go/gorilla/sdk/client/products"
)

func TestOutClient(t *testing.T) {
	config := client.DefaultTransportConfig().WithHost("localhost:9100")

	c := client.NewHTTPClientWithConfig(nil, config)
	params := products.NewListProductsParams()

	product, err := c.Products.ListProducts(params)
	fmt.Printf("%#v", product.GetPayload()[0])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(product)

}
