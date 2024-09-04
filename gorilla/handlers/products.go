package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/SergiuPlesco/microservices-go/gorilla/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, _ *http.Request) {
	pl := data.GetProducts()
	err := pl.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, errId := strconv.Atoi(vars["id"])
	if errId != nil {
		http.Error(rw, "Unable to convert id to number", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	errUpdate := data.UpdateProduct(id, prod)
	if errUpdate == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if errUpdate != nil {
		http.Error(rw, "Internal error", http.StatusInternalServerError)
		return
	}

}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to parse json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
