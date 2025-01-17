package handler_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	api_buyer  = "/api/v1/buyers/report-purchase-orders"
	txdb_buyer = "txdb_buyer"
)

type BuyerTestSuite struct {
	db *sql.DB
	hd *handler.BuyerHandlerDefault
	suite.Suite
}

func init() {
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "meli_pass",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "melifresh_db_buyer_test",
		ParseTime: true,
	}
	txdb.Register(txdb_buyer, "mysql", cfg.FormatDSN())
}

func (b *BuyerTestSuite) SetupTest() {
	var err error
	b.db, err = sql.Open(txdb_buyer, "")
	require.NoError(b.T(), err)
	rpBuyer := repository.NewBuyerMysqlRepository(b.db)
	svBuyer := service.NewBuyerService(rpBuyer)
	b.hd = handler.NewBuyerHandlerDefault(svBuyer)
}

func TestBuyerTestSuite(t *testing.T) {
	suite.Run(t, new(BuyerTestSuite))
}

func (b *BuyerTestSuite) TestReportPurchaseOrders() {
	defer b.db.Close()

	b.T().Run("case 1 - success - Purchase order reports for all buyers", func(t *testing.T) {
		// given
		request := httptest.NewRequest(http.MethodGet, api_buyer, nil)
		response := httptest.NewRecorder()

		// when
		b.hd.ReportPurchaseOrders(response, request)

		// then
		expectedCode := http.StatusOK
		expectedBody := `{
			"data": [
				{
					"id": 1,
					"card_number_id": "B1001",
					"first_name": "Alice",
					"last_name": "Brown",
					"purchase_orders_count": 1
				},
				{
					"id": 2,
					"card_number_id": "B1002",
					"first_name": "Mark",
					"last_name": "Jones",
					"purchase_orders_count": 1
				}, 
				{
					"id": 3,
					"card_number_id": "B1003",
					"first_name": "Linda",
					"last_name": "Garcia",
					"purchase_orders_count": 0
				}
			]
		}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	b.T().Run("case 2 - success - Purchase order reports by buyer id", func(t *testing.T) {
		// given
		request := httptest.NewRequest(http.MethodGet, api_buyer+"?id=2", nil)
		response := httptest.NewRecorder()

		// when
		b.hd.ReportPurchaseOrders(response, request)

		// then
		expectedCode := http.StatusOK
		expectedBody := `{
			"data": [
				{
					"id": 2,
					"card_number_id": "B1002",
					"first_name": "Mark",
					"last_name": "Jones",
					"purchase_orders_count": 1
				}
			]
		}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	b.T().Run("case 3 - success - Purchase order reports by buyer id that doesn't have purchase orders", func(t *testing.T) {
		// given
		request := httptest.NewRequest(http.MethodGet, api_buyer+"?id=3", nil)
		response := httptest.NewRecorder()

		// when
		b.hd.ReportPurchaseOrders(response, request)

		// then
		expectedCode := http.StatusOK
		expectedBody := `{
			"data": [
				{
					"id": 3,
					"card_number_id": "B1003",
					"first_name": "Linda",
					"last_name": "Garcia",
					"purchase_orders_count": 0
				}
			]
		}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	b.T().Run("case 4 - error - Purchase order reports by buyer id that doesn't exists", func(t *testing.T) {
		// given
		request := httptest.NewRequest(http.MethodGet, api_buyer+"?id=4", nil)
		response := httptest.NewRecorder()

		// when
		b.hd.ReportPurchaseOrders(response, request)

		// then
		expectedCode := http.StatusNotFound
		expectedBody := `{
			"message": "buyer not found",
			"error": "not_found",
			"code": 404,
			"causes": null
		}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})
}
