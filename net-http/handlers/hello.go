package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte("Ooops"))
		return
	}
	fmt.Fprintf(rw, "%s\n", data)
}
