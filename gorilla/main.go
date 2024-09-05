// Package classification of Product API
//
// Documentation for Product API
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SergiuPlesco/microservices-go/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "gorilla", log.LstdFlags)

	// create the handlers
	pl := handlers.NewProducts(l)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", pl.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", pl.UpdateProduct)
	putRouter.Use(pl.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", pl.AddProduct)
	postRouter.Use(pl.MiddlewareProductValidation)

	// create new server
	s := &http.Server{
		Addr:         ":9100",           // configure the bind address
		Handler:      sm,                //set the default handler
		ErrorLog:     l,                 // set the logger for the server
		IdleTimeout:  120 * time.Second, // max time to read request from the client
		ReadTimeout:  1 * time.Second,   // max time to write response to the client
		WriteTimeout: 1 * time.Second,   // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		log.Println("Starting server on :9100...")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	// trap sigterm or interrupt and gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	l.Println("Recieved tarminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	s.Shutdown(tc)

	log.Println("Server exited properly")
}
