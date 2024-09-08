package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SergiuPlesco/microservices-go/gorilla/data"
	"github.com/SergiuPlesco/microservices-go/gorilla/handlers"
	"github.com/SergiuPlesco/microservices-go/grpc/pb"
	"github.com/go-openapi/runtime/middleware"
	corsHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	l := hclog.Default()
	v := data.NewValidation()

	conn, err := grpc.NewClient("localhost:9101", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	cc := pb.NewCurrencyClient(conn)

	db := data.NewProductsDB(cc, l)

	// create the handlers
	pl := handlers.NewProducts(l, v, db)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", pl.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", pl.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", pl.Update)
	putRouter.Use(pl.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", pl.Create)
	postRouter.Use(pl.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("products/{id:[0-9]+}", pl.Delete)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	ch := corsHandlers.CORS(corsHandlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost:3001"}))

	// create new server
	s := &http.Server{
		Addr:         ":9100",                                          // configure the bind address
		Handler:      ch(sm),                                           //set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		IdleTimeout:  120 * time.Second,                                // max time to read request from the client
		ReadTimeout:  1 * time.Second,                                  // max time to write response to the client
		WriteTimeout: 1 * time.Second,                                  // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		log.Println("Starting server on :9100...")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server", "error", err)
		}
	}()
	// trap sigterm or interrupt and gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	l.Info("Recieved tarminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	s.Shutdown(tc)

	log.Println("Server exited properly")
}
