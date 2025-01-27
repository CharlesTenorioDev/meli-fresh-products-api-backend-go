//go:build integration

package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const inboundRoute = "/api/v1/inbound-orders"

type InboundOrdersTestSuite struct {
	hd *handler.InboundOrdersHandler
	db *sql.DB
	suite.Suite
}

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "meli_pass",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "melifresh",
	}
	txdb.Register("txdb_inboundRoute", "mysql", cfg.FormatDSN())
}

func (suite *InboundOrdersTestSuite) SetupTest() {
	var err error
	suite.db, err = sql.Open("txdb_inboundRoute", "identier")
	require.NoError(suite.T(), err)

	tempFile, err := os.CreateTemp("", "warehouse.json")
	require.NoError(suite.T(), err)
	defer os.Remove(tempFile.Name()) // clean up

	warehouseJSON := `[{"warehouse_code": "WHS001", "address": "Rua Exemplo, 123", "telephone": "123456789", "minimum_capacity": 100, "minimum_temperature": 5.0}]`
	err = os.WriteFile(tempFile.Name(), []byte(warehouseJSON), 0644)
	require.NoError(suite.T(), err)

	rpInboundOrder := repository.NewInboundOrderMysql(suite.db)
	rpEmployee := repository.NewEmployeeMysql(suite.db)
	rpProductBatch := repository.NewProductBatchMysql(suite.db)
	rpWarehouse := repository.NewRepositoryWarehouse(nil, tempFile.Name())
	sv := service.NewInboundOrderService(rpInboundOrder, rpEmployee, rpProductBatch, rpWarehouse)
	suite.hd = handler.NewInboundOrdersHandler(sv)
}

func (suite *InboundOrdersTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *InboundOrdersTestSuite) TestInboundOrdersHandler_Create() {
	defer suite.db.Close()

	suite.T().Run("success to create with valid data", func(t *testing.T) {

		inbound := internal.InboundOrders{
			OrderDate:      "2025-01-01",
			OrderNumber:    "ORD5555",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		b, err := json.Marshal(inbound)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, inboundRoute, bytes.NewReader(b))
		w := httptest.NewRecorder()

		suite.hd.Create(w, r)

		var res struct {
			Data struct {
				Id int64 `json:"id"`
			} `json:"data"`
		}
		json.NewDecoder(w.Result().Body).Decode(&res)
		require.Equal(t, http.StatusCreated, w.Result().StatusCode)
		require.Greater(t, res.Data.Id, int64(0))
	})

	suite.T().Run("attempt to create with missing required fields", func(t *testing.T) {

		inbound := internal.InboundOrders{}
		b, err := json.Marshal(inbound)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, inboundRoute, bytes.NewReader(b))
		w := httptest.NewRecorder()

		suite.hd.Create(w, r)

		require.Equal(t, http.StatusUnprocessableEntity, w.Result().StatusCode)
	})

	suite.T().Run("conflit to create with existing order number", func(t *testing.T) {

		inbound := internal.InboundOrders{
			OrderDate:      "2025-01-01",
			OrderNumber:    "ORD001",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		b, err := json.Marshal(inbound)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, inboundRoute, bytes.NewReader(b))
		w := httptest.NewRecorder()

		suite.hd.Create(w, r)
		require.Equal(t, http.StatusConflict, w.Result().StatusCode)
	})

	suite.T().Run("create with not exists foreign key reference (employee_id)", func(t *testing.T) {

		inbound := internal.InboundOrders{
			OrderDate:      "2025-01-01",
			OrderNumber:    "ORD002",
			EmployeeID:     999,
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		b, err := json.Marshal(inbound)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, inboundRoute, bytes.NewReader(b))
		w := httptest.NewRecorder()

		suite.hd.Create(w, r)
		require.Equal(t, http.StatusConflict, w.Result().StatusCode)
	})
}

func (suite *InboundOrdersTestSuite) TestInboundOrdersHandler_GetAll() {
	defer suite.db.Close()

	suite.db.Exec("INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES ('2025-01-01', 'ORD001', 1, 1, 1)")

	r := httptest.NewRequest(http.MethodGet, inboundRoute, nil)
	w := httptest.NewRecorder()

	suite.hd.GetAll(w, r)

	var res struct {
		Data []internal.InboundOrders `json:"data"`
	}
	json.NewDecoder(w.Result().Body).Decode(&res)

	require.Greater(suite.T(), len(res.Data), 0)
}

func TestInboundOrdersTestSuite(t *testing.T) {
	suite.Run(t, new(InboundOrdersTestSuite))
}
