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

	whRepository := repository.NewRepositoryWarehouse(nil, "db/warehouse.json")
	slRepository := repository.NewSellerRepoMap()

	rt.Route("/api/v1", func(r chi.Router) {
		r.Route("/employees", employeeRouter)
		r.Route("/buyers", buyerRouter)
		r.Route("/sections", func(r chi.Router) {
			sectionsRoutes(r, whRepository)
		})
		r.Route("/warehouses", func(r chi.Router) {
			warehouseRoute(r, whRepository)
		})
		r.Route("/sellers", func(r chi.Router) {
			sellerRoutes(r, slRepository)
		})
		r.Route("/products", func(r chi.Router) {
			productRoutes(r, slRepository)
		})
	})

	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func sellerRoutes(r chi.Router, slRepository internal.SellerRepository) {
	sv := service.NewSellerServiceDefault(slRepository)
	hd := handler.NewSellerDefault(sv)

	r.Get("/", hd.GetAll())
	r.Get("/{id}", hd.GetByID())
	r.Post("/", hd.Save())
	r.Patch("/{id}", hd.Update())
	r.Delete("/{id}", hd.Delete())
}

func warehouseRoute(r chi.Router, whRepository internal.WarehouseRepository) {
	warehouseService := service.NewWarehouseDefault(whRepository)
	warehouseHandler := handler.NewWarehouseDefault(warehouseService)

	r.Get("/", warehouseHandler.GetAll())
	r.Get("/{id}", warehouseHandler.GetByID())
	r.Post("/", warehouseHandler.Create())
	r.Patch("/{id}", warehouseHandler.Update())
	r.Delete("/{id}", warehouseHandler.Delete())
}

func sectionsRoutes(r chi.Router, whRepository internal.WarehouseRepository) {
	rpS := repository.NewRepositorySection()
	rpT := repository.NewRepositoryProductType()
	sv := service.NewServiceSection(rpS, rpT, whRepository)
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

func productRoutes(r chi.Router, slRepository internal.SellerRepository) {
	repo := repository.NewProductMap()
	svc := service.NewProductService(repo, slRepository)
	hd := handler.NewProducHandlerDefault(svc)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}",hd.Update)
	r.Delete("/{id}", hd.Delete) 
}