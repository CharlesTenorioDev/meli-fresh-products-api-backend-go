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
	router := chi.NewRouter()

	rpS := repository.NewRepositorySection()
	rpW := repository.NewRepositoryWareHouse()

	sv := service.NewServiceSection(rpS, rpW)
	hd := handler.NewHandlerSection(sv)

	router.Route("/api/v1/sections", func(r chi.Router) {
		r.Get("/", hd.GetAll)
		r.Get("/{id}", hd.GetByID)
		r.Post("/", hd.Create)
		r.Patch("/{id}", hd.Update)
		r.Delete("/{id}", hd.Delete)
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
