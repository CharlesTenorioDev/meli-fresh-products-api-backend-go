package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
)

type ConfigServerChi struct {
	ServerAddress string
}

func NewServerChi(cfg *ConfigServerChi) *ServerChi {
	defaultConfig := &ConfigServerChi{
		ServerAddress: ":8080",
	}
	if cfg != nil {
		if cfg.ServerAddress != "" {
			defaultConfig.ServerAddress = cfg.ServerAddress
		}

	}

	return &ServerChi{
		serverAddress: defaultConfig.ServerAddress,
	}
}

type ServerChi struct {
	serverAddress string
}

func (a *ServerChi) Run() (err error) {
	rt := chi.NewRouter()
	rt.Use(middleware.Logger)

	

	rt.Route("/api/v1", func(r chi.Router) {
		r.Route("/products", ProductRouter)
	})

	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func ProductRouter(r chi.Router) {
	repo := repository.NewProductMap()
	svc := service.NewProductService(repo)
	hd := handler.NewProducHandlerDefault(svc)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
}