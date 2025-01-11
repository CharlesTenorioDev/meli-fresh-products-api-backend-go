package application

import (
	"database/sql"
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
	Dsn           string
}

// NewServerChi is a function that returns a new instance of ServerChi
func NewServerChi(cfg *ConfigServerChi) *ServerChi {
	// default values
	defaultConfig := &ConfigServerChi{
		ServerAddress: ":8080",
		Dsn:           "",
	}
	if cfg != nil {
		if cfg.ServerAddress != "" {
			defaultConfig.ServerAddress = cfg.ServerAddress
		}
		if cfg.Dsn != "" {
			defaultConfig.Dsn = cfg.Dsn
		}

	}

	return &ServerChi{
		serverAddress: defaultConfig.ServerAddress,
		dsn:           defaultConfig.Dsn,
	}
}

// ServerChi is a struct that implements the Application interface
type ServerChi struct {
	// serverAddress is the address where the server will be listening
	serverAddress string
	dsn           string
}

// Run is a method that runs the application
func (a *ServerChi) Run() (err error) {

	db, err := sql.Open("mysql", a.dsn)
	if err != nil {
		return
	}

	defer db.Close()

	// - database: ping
	err = db.Ping()
	if err != nil {
		return
	}

	rt := chi.NewRouter()
	rt.Use(middleware.Logger)

	whRepository := repository.NewRepositoryWarehouse(nil, "db/warehouse.json")
	slRepository := repository.NewSellerMysql(db)
	lcRepository := repository.NewLocalityMysql(db)
	pdRepository := repository.NewProductMap()
	empRepository := repository.NewEmployeeMysql(db)
	inbRepository := repository.NewInboundOrderMysql(db)

	rt.Route("/api/v1", func(r chi.Router) {
		r.Route("/employees", func(r chi.Router) {
			employeeRouter(r, whRepository, db)
		})
		r.Route("/buyers", buyerRouter)
		r.Route("/sections", func(r chi.Router) {
			sectionsRoutes(r, whRepository, pdRepository)
		})
		r.Route("/warehouses", func(r chi.Router) {
			warehouseRoute(r, whRepository)
		})
		r.Route("/sellers", func(r chi.Router) {
			sellerRoutes(r, slRepository, lcRepository)
		})
		r.Route("/localities", func(r chi.Router) {
			localitiesRoutes(r, lcRepository)
		})
		r.Route("/products", func(r chi.Router) {
			productRoutes(r, pdRepository, slRepository)
		})
		r.Route("/inbound-orders", func(r chi.Router) {
			inboundOrdersRoutes(r, inbRepository, empRepository, whRepository)
		})
	})

	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func localitiesRoutes(r chi.Router, lcRepository internal.LocalityRepository) {
	sv := service.NewLocalityDefault(lcRepository)
	hd := handler.NewLocalityDefault(sv)

	r.Get("/report-sellers", hd.ReportSellers())
	r.Post("/", hd.Save())
}

func sellerRoutes(r chi.Router, slRepository internal.SellerRepository, lcRepository internal.LocalityRepository) {
	sv := service.NewSellerServiceDefault(slRepository, lcRepository)
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

func sectionsRoutes(r chi.Router, whRepository internal.WarehouseRepository, ptRepository internal.ProductRepository) {
	rpS := repository.NewRepositorySection()
	rpT := repository.NewRepositoryProductType()
	sv := service.NewServiceSection(rpS, rpT, ptRepository, whRepository)
	hd := handler.NewHandlerSection(sv)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
}

func employeeRouter(r chi.Router, whRepository internal.WarehouseRepository, db *sql.DB) {
	rp := repository.NewEmployeeMysql(db)
	sv := service.NewEmployeeServiceDefault(rp, whRepository)
	hd := handler.NewEmployeeDefault(sv)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
	r.Get("/report-inbound-orders", hd.ReportInboundOrders)
}

func buyerRouter(r chi.Router) {
	repo := repository.NewBuyerMap("db/buyer.json")
	svc := service.NewBuyerService(repo)
	hd := handler.NewBuyerHandlerDefault(svc)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
}

func productRoutes(r chi.Router, ptRepo internal.ProductRepository, slRepository internal.SellerRepository) {
	rpT := repository.NewRepositoryProductType()
	svc := service.NewProductService(ptRepo, slRepository, rpT)
	hd := handler.NewProducHandlerDefault(svc)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
}

func inboundOrdersRoutes(r chi.Router, inbRepository internal.InboundOrdersRepository, empRepository internal.EmployeeRepository, whRepository internal.WarehouseRepository) {
	sv := service.NewInboundOrderService(inbRepository, empRepository, whRepository)
	hd := handler.NewInboundOrdersHandler(sv)

	r.Post("/", hd.Create)
	r.Get("/", hd.GetAll)

}
