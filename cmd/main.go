package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cr := chi.NewRouter()

	cr.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	})
	if err := http.ListenAndServe(":8080", cr); err != nil {
		log.Fatal(err)
	}
}
