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

	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	// - endpoints
	employeeRoutes(rt)

	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func employeeRoutes(router *chi.Mux) {
	// - repository
	rp := repository.NewEmployeeRepository()
	// - service
	sv := service.NewEmployeeServiceDefault(rp)
	// - handler
	hd := handler.NewEmployeeDefault(sv)

	router.Route("/api/v1/employees", func(rt chi.Router) {
		rt.Get("/", hd.GetAll)
		rt.Get("/{id}", hd.GetByID)
		// rt.Post("/", hd.Save)
		// rt.Patch("/{id}", hd.Update)
		// rt.Delete("/{id}", hd.Delete)
	})
}
