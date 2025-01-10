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
	scRepository := repository.NewRepositorySectionMysql(db)
	pbRepository := repository.NewRepositoryProductBatchMysql(db)

	rt.Route("/api/v1", func(r chi.Router) {
		r.Route("/employees", func(r chi.Router) {
			employeeRouter(r, whRepository)
		})
		r.Route("/buyers", buyerRouter)
		r.Route("/sections", func(r chi.Router) {
			sectionsRoutes(r, scRepository, whRepository, pdRepository)
		})
		r.Route("/product-batches", func(r chi.Router) {
			productBatchRoutes(r, pbRepository, scRepository, pdRepository)
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
		r.Route("/carries", func(r chi.Router) {
			carriesRoutes(r, db)
		})
	})

	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func localitiesRoutes(r chi.Router, lcRepository internal.LocalityRepository) {
	sv := service.NewLocalityDefault(lcRepository)
	hd := handler.NewLocalityDefault(sv)

	r.Get("/report-sellers", hd.ReportSellers())
	r.Get("/report-carries", hd.ReportCarries())
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

func sectionsRoutes(r chi.Router, scRepository internal.SectionRepository, whRepository internal.WarehouseRepository, ptRepository internal.ProductRepository) {
	rpT := repository.NewRepositoryProductType()
	sv := service.NewServiceSection(scRepository, rpT, ptRepository, whRepository)
	hd := handler.NewHandlerSection(sv)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Get("/report-products", hd.ReportProducts)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
}

func productBatchRoutes(r chi.Router, pbRepository internal.ProductBatchRepository, scRepository internal.SectionRepository, ptRepository internal.ProductRepository) {
	sv := service.NewServiceProductBatch(pbRepository, scRepository, ptRepository)
	hd := handler.NewHandlerProductBatch(sv)

	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
}

func employeeRouter(r chi.Router, whRepository internal.WarehouseRepository) {
	rp := repository.NewEmployeeRepository()
	sv := service.NewEmployeeServiceDefault(rp, whRepository)
	hd := handler.NewEmployeeDefault(sv)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
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

func carriesRoutes(r chi.Router, db *sql.DB) {
	rp := repository.NewCarriesMysql(db)
	sv := service.NewCarriesService(rp)
	hd := handler.NewCarriesHandlerDefault(sv)

	r.Get("/", hd.GetAll)
	r.Post("/", hd.Create)
}
