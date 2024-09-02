package main

import (
	"log"
	"net/http"
)

func main() {

	log.Println("Starting server on :9090...")
	if err := http.ListenAndServe(":9100", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
