package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
)

func main() {
	cr := chi.NewRouter()
	repo := repository.NewBuyerMap()
	svc := service.NewBuyerService(repo)
	hd := handler.NewBuyerHandlerDefault(svc)

	cr.Route("/api/v1", func(r chi.Router) {
		r.Route("/buyers", func(r chi.Router) {
			r.Get("/", hd.GetAll())
		})
	})
	if err := http.ListenAndServe(":8080", cr); err != nil {
		log.Fatal(err)
	}
}
