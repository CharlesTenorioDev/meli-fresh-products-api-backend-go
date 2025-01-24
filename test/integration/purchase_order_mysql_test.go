//go:build integration

package integration_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	api_po    = "/api/v1/purchase-orders"
	txdb_name = "txdb_purchase_order"
)

type PurchaseOrdersTestSuite struct {
	db *sql.DB
	hd *handler.PurchaseOrderHandler
	rp *repository.PurchaseOrderRepository
	suite.Suite
}

func init() {
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "meli_pass",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "melifresh_purchase_orders_test_db",
		ParseTime: true,
	}
	txdb.Register(txdb_name, "mysql", cfg.FormatDSN())
}

func (p *PurchaseOrdersTestSuite) SetupTest() {
	// dependencies
	var err error
	p.db, err = sql.Open(txdb_name, "melifresh_purchase_orders_test_db")
	require.NoError(p.T(), err)
	rpPurchaseOrder := repository.NewPurchaseOrderMysqlRepository(p.db)
	rpBuyer := repository.NewBuyerMysqlRepository(p.db)
	rpProductRecord := repository.NewProductRecordsSQL(p.db)
	svBuyer := service.NewBuyerService(rpBuyer)
	sv := service.NewPurchaseOrderService(rpPurchaseOrder, rpProductRecord, svBuyer)

	p.rp = rpPurchaseOrder
	p.hd = handler.NewPurchaseOrderHandler(sv)
}

func (p *PurchaseOrdersTestSuite) TestPurchaseOrders_Create() {
	defer p.db.Close()
	p.T().Run("case 1: success - creating a purchase order with valid data", func(t *testing.T) {
		// given
		request := &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Path: api_po},
			Body: io.NopCloser(strings.NewReader(
				`{
					"order_number": "123ABC",
					"order_date": "2023-10-05",
					"tracking_code": "123",
					"buyer_id": 1,
					"product_record_id": 2
				}`,
			)),
			Header: http.Header{"Content-Type": []string{"application/json"}},
		}
		response := httptest.NewRecorder()

		// when
		p.hd.Create()(response, request)
		var data struct {
			PurchaseOrderCreated internal.PurchaseOrder `json:"data"`
		}
		err := json.NewDecoder(response.Result().Body).Decode(&data)
		if err != nil {
			t.Fatal(err)
		}
		// retrieve the purchase order from the database
		var po internal.PurchaseOrder
		po, errPo := p.rp.FindByID(data.PurchaseOrderCreated.ID)
		require.NoError(t, errPo)

		// then
		expectedCode := http.StatusCreated
		expectedBody := fmt.Sprintf(`{"data":{"id":%d,"order_number":"123ABC","order_date":"2023-10-05","tracking_code":"123","buyer_id":1,"product_record_id":2}}`, po.ID)
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeader, response.Header())
	})

	p.T().Run("case 2: error - creating a purchase order with a duplicated order number", func(t *testing.T) {
		// given
		request := &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Path: api_po},
			Body: io.NopCloser(strings.NewReader(
				`{
					"order_number": "PO1001",
					"order_date": "2023-10-05",
					"tracking_code": "123",
					"buyer_id": 1,
					"product_record_id": 2
				}`,
			)),
			Header: http.Header{"Content-Type": []string{"application/json"}},
		}
		response := httptest.NewRecorder()

		// when
		p.hd.Create()(response, request)

		// then
		expectedCode := http.StatusConflict
		expectedBody := `{
			"message": "purchase order number already exists",
			"error": "conflict",
			"code": 409,
			"causes": null
		}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeader, response.Header())
	})

	p.T().Run("case 3: error - creating a purchase order with an invalid date format", func(t *testing.T) {
		// given
		request := &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Path: api_po},
			Body: io.NopCloser(strings.NewReader(
				`{
					"order_number": "PO1001",
					"order_date": "2023-10-05asdf",
					"tracking_code": "123",
					"buyer_id": 1,
					"product_record_id": 2
				}`,
			)),
			Header: http.Header{"Content-Type": []string{"application/json"}},
		}
		response := httptest.NewRecorder()

		// when
		p.hd.Create()(response, request)

		// then
		expectedCode := http.StatusBadRequest
		expectedBody := `{
			"message": "Invalid data",
			"error": "bad_request",
			"code": 400,
			"causes": [
				{
					"field": "order_date",
					"message": "invalid date format"
				}
			]
		}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeader, response.Header())
	})

	p.T().Run("case 4: error - creating a purchase order with an invalid buyer id", func(t *testing.T) {
		// given
		request := &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Path: api_po},
			Body: io.NopCloser(strings.NewReader(
				`{
					"order_number": "PO1001",
					"order_date": "2023-10-05",
					"tracking_code": "123",
					"buyer_id": 100,
					"product_record_id": 2
				}`,
			)),
			Header: http.Header{"Content-Type": []string{"application/json"}},
		}
		response := httptest.NewRecorder()

		// when
		p.hd.Create()(response, request)

		// then
		expectedCode := http.StatusNotFound
		expectedBody := `{
			"message": "buyer not found",
			"error": "not_found",
			"code": 404,
			"causes": null
		}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeader, response.Header())
	})

	p.T().Run("case 5: error - creating a purchase order with an invalid product record id", func(t *testing.T) {
		// given
		request := &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Path: api_po},
			Body: io.NopCloser(strings.NewReader(
				`{
					"order_number": "PO1001",
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
		expectedCode := http.StatusNotFound
		expectedBody := `{
			"message": "product-records not found",
			"error": "not_found",
			"code": 404,
			"causes": null
		}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeader, response.Header())
	})

	p.T().Run("case 6: error - creating a purchase order with empty inputs", func(t *testing.T) {
		// given
		request := &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Path: api_po},
			Body: io.NopCloser(strings.NewReader(`{
				"order_date": "2023-10-05"
			}`)),
			Header: http.Header{"Content-Type": []string{"application/json"}},
		}
		response := httptest.NewRecorder()

		// when
		p.hd.Create()(response, request)

		// then
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{
			"message": "Purchase Order inputs are Invalid",
			"error": "bad_request",
			"code": 400,
			"causes": [
				{
					"field": "order_number",
					"message": "Order number is required"
				},
				{
					"field": "tracking_code",
					"message": "Tracking code is required"
				},
				{
					"field": "buyer_id",
					"message": "Buyer ID is required"
				},
				{
					"field": "product_record_id",
					"message": "Product record ID is required"
				}
			]
		}`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
		require.Equal(t, expectedHeader, response.Header())
	})
}

func TestPurchaseOrdersTestSuite(t *testing.T) {
	suite.Run(t, new(PurchaseOrdersTestSuite))
}
