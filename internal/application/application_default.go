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
	_ "github.com/meli-fresh-products-api-backend-t1/swagger/docs"
	httpSwagger "github.com/swaggo/http-swagger"
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
		return err
	}

	defer db.Close()

	// - database: ping
	err = db.Ping()
	if err != nil {
		return err
	}

	rt := chi.NewRouter()
	rt.Use(middleware.Logger)
	rt.Get("/swagger/*", httpSwagger.WrapHandler)

	buMysqlRepository := repository.NewBuyerMysqlRepository(db)
	whRepository := repository.NewWarehouseMysqlRepository(db)
	slRepository := repository.NewSellerMysql(db)
	lcRepository := repository.NewLocalityMysql(db)
	pdRepository := repository.NewProductSQL(db)
	prodRecRepository := repository.NewProductRecordsSQL(db)
	emRepository := repository.NewEmployeeMysql(db)
	inRepository := repository.NewInboundOrderMysql(db)
	scRepository := repository.NewSectionMysql(db)
	pbRepository := repository.NewProductBatchMysql(db)
	ptRepository := repository.NewProductTypeMysql(db)
	poMysqlRepository := repository.NewPurchaseOrderMysqlRepository(db)
	buyerService := service.NewBuyerService(buMysqlRepository)

	rt.Route("/api/v1", func(r chi.Router) {
		r.Route("/employees", func(r chi.Router) {
			employeeRouter(r, whRepository, db)
		})
		r.Route("/buyers", func(r chi.Router) {
			buyerRouter(r, buMysqlRepository)
		})
		r.Route("/sections", func(r chi.Router) {
			sectionsRoutes(r, scRepository, ptRepository, whRepository, pdRepository)
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
			productRoutes(r, pdRepository, slRepository, ptRepository)
		})
		r.Route("/purchase-orders", func(r chi.Router) {
			purchaseOrderRouter(r, poMysqlRepository, prodRecRepository, buyerService)
		})
		r.Route("/carries", func(r chi.Router) {
			carriesRoutes(r, db)
		})

		r.Route("/productRecords", func(r chi.Router) {
			productRecordsRoutes(r, prodRecRepository, pdRepository)
		})

		r.Route("/inbound-orders", func(r chi.Router) {
			inboundOrdersRoutes(r, inRepository, emRepository, pbRepository, whRepository)
		})
	})

	err = http.ListenAndServe(a.serverAddress, rt)

	return err
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

func sectionsRoutes(r chi.Router, scRepository internal.SectionRepository, ptRepository internal.ProductTypeRepository, whRepository internal.WarehouseRepository, pdRepository internal.ProductRepository) {
	sv := service.NewServiceSection(scRepository, ptRepository, pdRepository, whRepository)
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

func buyerRouter(r chi.Router, buRepository internal.BuyerRepository) {
	svc := service.NewBuyerService(buRepository)
	hd := handler.NewBuyerHandlerDefault(svc)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
	r.Get("/report-purchase-orders", hd.ReportPurchaseOrders)
}

func productRoutes(r chi.Router, pdRepository internal.ProductRepository, slRepository internal.SellerRepository, ptRepository internal.ProductTypeRepository) {
	svc := service.NewProductService(pdRepository, slRepository, ptRepository)
	hd := handler.NewProductHandlerDefault(svc)

	r.Get("/", hd.GetAll)
	r.Get("/{id}", hd.GetByID)
	r.Post("/", hd.Create)
	r.Patch("/{id}", hd.Update)
	r.Delete("/{id}", hd.Delete)
	r.Get("/report-records", hd.ReportRecords)
}

func inboundOrdersRoutes(r chi.Router, inRepository internal.InboundOrdersRepository, emRepository internal.EmployeeRepository, pbRepository internal.ProductBatchRepository, whRepository internal.WarehouseRepository) {
	sv := service.NewInboundOrderService(inRepository, emRepository, pbRepository, whRepository)
	hd := handler.NewInboundOrdersHandler(sv)

	r.Post("/", hd.Create)
	r.Get("/", hd.GetAll)
}

func purchaseOrderRouter(r chi.Router, poRepository internal.PurchaseOrderRepository, prodRecRepository internal.ProductRecordsRepository, buyerService internal.BuyerService) {
	sv := service.NewPurchaseOrderService(poRepository, prodRecRepository, buyerService)
	hd := handler.NewPurchaseOrderHandler(sv)

	r.Post("/", hd.Create())
}

func carriesRoutes(r chi.Router, db *sql.DB) {
	rp := repository.NewCarriesMysql(db)
	sv := service.NewCarriesService(rp)
	hd := handler.NewCarriesHandlerDefault(sv)

	r.Get("/", hd.GetAll)
	r.Post("/", hd.Create)
}
func productRecordsRoutes(r chi.Router, prodRecRepository internal.ProductRecordsRepository, prodRepository internal.ProductRepository) {
	svc := service.NewProductRecordsDefault(prodRecRepository, prodRepository)
	hd := handler.NewProductRecordsDefault(svc)
	r.Post("/", hd.Create)
}
