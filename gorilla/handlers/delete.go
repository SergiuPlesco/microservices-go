package handlers

import (
	"net/http"

	"github.com/SergiuPlesco/microservices-go/gorilla/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)
	p.l.Debug("[DEBUG] deleting record id", id)

	err := p.productDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("[ERROR] deleting record id does not exists")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)

	}

}
