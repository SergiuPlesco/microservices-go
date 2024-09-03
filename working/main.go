package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SergiuPlesco/microservices-go/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	sm := http.NewServeMux()

	sm.Handle("/", hh)
	sm.Handle("goodbye", gh)

	log.Println("Starting server on :9090...")
	if err := http.ListenAndServe(":9100", sm); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
