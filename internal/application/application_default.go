package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"net/http"
)

// ConfigServerChi is a struct that represents the configuration for ServerChi
type ConfigServerChi struct {
	// ServerAddress is the address where the server will be listening
	ServerAddress string
}

// NewServerChi is a function that returns a new instance of ServerChi
func NewServerChi(cfg *ConfigServerChi) *ServerChi {
	// default values
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

// ServerChi is a struct that implements the Application interface
type ServerChi struct {
	// serverAddress is the address where the server will be listening
	serverAddress string
}

// Run is a method that runs the application
func (a *ServerChi) Run() (err error) {
	// dependencies
	// - loader

	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	// - endpoints
	sellerRoutes(rt)

	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func sellerRoutes(router *chi.Mux) {
	// - repository
	rp := repository.NewSellerRepoMap(make(map[int]internal.Seller))
	// - service
	sv := service.NewSellerServiceDefault(rp)
	// - handler
	hd := handler.NewSellerDefault(sv)

	router.Route("/api/v1/sellers", func(rt chi.Router) {
		rt.Get("/", hd.GetAll())
		rt.Get("/{id}", hd.GetByID())
		rt.Post("/", hd.Save())
		rt.Patch("/{id}", hd.Update())
		rt.Delete("/{id}", hd.Delete())
	})
}
