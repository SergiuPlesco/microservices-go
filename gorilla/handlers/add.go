package handlers

import (
	"net/http"

	"github.com/SergiuPlesco/microservices-go/gorilla/data"
)

// swagger:route GET /products products product list
// Returns a list of products
// responses:
// 200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, _ *http.Request) {
	pl := data.GetProducts()
	err := pl.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
