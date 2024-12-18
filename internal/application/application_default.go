package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
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
	rt := chi.NewRouter()
	rt.Use(middleware.Logger)

	rt.Route("/api/v1", func(r chi.Router) {
		r.Route("/sellers", sellerRoutes)
		r.Route("/warehouses", warehouseRoute)
		r.Route("/sections", sectionsRoutes)
		r.Route("/employees", employeeRouter)
		r.Route("/buyers", buyerRouter)
		r.Route("/products", productRouter)
	})

	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func sellerRoutes(r chi.Router) {
	rp := repository.NewSellerRepoMap(make(map[int]internal.Seller))
	sv := service.NewSellerServiceDefault(rp)
	hd := handler.NewSellerDefault(sv)

	r.Get("/", hd.GetAll())
	r.Get("/{id}", hd.GetByID())
	r.Post("/", hd.Save())
	r.Patch("/{id}", hd.Update())
	r.Delete("/{id}", hd.Delete())
}

func warehouseRoute(r chi.Router) {
	warehouseRepository := repository.NewRepositoryWarehouse(nil, "db/warehouse.json")
	warehouseService := service.NewWarehouseDefault(warehouseRepository)
	warehouseHandler := handler.NewWarehouseDefault(warehouseService)

	r.Get("/", warehouseHandler.GetAll())
	r.Get("/{id}", warehouseHandler.GetByID())
	r.Post("/", warehouseHandler.Create())
	r.Patch("/{id}", warehouseHandler.Update())
	r.Delete("/{id}", warehouseHandler.Delete())
}

func sectionsRoutes(r chi.Router) {
	rpS := repository.NewRepositorySection()
	rpW := repository.NewRepositoryWarehouse(nil, "db/warehouse.json")

	sv := service.NewServiceSection(rpS, rpW)
	hd := handler.NewHandlerSection(sv)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
}

func employeeRouter(r chi.Router) {
	repo := repository.NewEmployeeRepository()
	svc := service.NewEmployeeServiceDefault(repo)
	hd := handler.NewEmployeeDefault(svc)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Save)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
}

func buyerRouter(r chi.Router) {
	repo := repository.NewBuyerMap()
	svc := service.NewBuyerService(repo)
	hd := handler.NewBuyerHandlerDefault(svc)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
}

func productRouter(r chi.Router) {
	repo := repository.NewProductMap()
	svc := service.NewProductService(repo)
	hd := handler.NewProducHandlerDefault(svc)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}",hd.Update)
	r.Delete("/{id}", hd.Delete) 
}