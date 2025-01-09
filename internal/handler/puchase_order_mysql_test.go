package handler_test

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const api = "/api/v1/purchase-orders"

type PurchaseOrdersTestSuite struct {
	hd *handler.PurchaseOrderHandler
	rt chi.Router
	db *sql.DB
	suite.Suite
}

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "meli_pass",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "melifresh_purchase_orders_test_db",
	}
	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func (p *PurchaseOrdersTestSuite) SetupTest() {
	var err error
	p.db, err = sql.Open("txdb", "")
	require.NoError(p.T(), err)
	rp := repository.NewPurchaseOrderMysqlRepository(p.db)
	rpBuyer := repository.NewBuyerMysqlRepository(p.db)
	svBuyer := service.NewBuyerService(rpBuyer)
	sv := service.NewPurchaseOrderService(rp, svBuyer)
	p.hd = handler.NewPurchaseOrderHandler(sv)
}

func (p *PurchaseOrdersTestSuite) TestPurchaseOrders_Create() {
	defer p.db.Close()

	p.T().Run("Creating a purchase order with valid data", func(t *testing.T) {
		// given
		request := &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Path: api},
			Body: io.NopCloser(strings.NewReader(
				`{
					"order_number": "123ABC",
					"order_date": "2023-10-05",   
					"tracking_code": "123",
					"buyer_id": 1,
					"product_record_id": 100
				}`,
			)),
			Header: http.Header{"Content-Type": []string{"application/json"}},
		}
		response := httptest.NewRecorder()

		// when
		p.hd.Create()(response, request)

		// then
		expectedCode := http.StatusCreated
		expectedBody := `{"data": {"id":11,"order_number":"123ABC","order_date":"2023-10-05T00:00:00Z","tracking_code":"123","buyer_id":1,"product_record_id":100}}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeader, response.Header())
	})
}

func TestPurchaseOrdersTestSuite(t *testing.T) {
	suite.Run(t, new(PurchaseOrdersTestSuite))
}
