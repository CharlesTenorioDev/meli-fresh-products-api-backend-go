package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	employeeRoute = "/api/v1/employees"
)

type EmployeeTestSuite struct {
	db *sql.DB
	rt chi.Router
	hd *handler.EmployeeHandlerDefault
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
	txdb.Register("txdb_employee", "mysql", cfg.FormatDSN())
}

type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeService) CountInboundOrdersPerEmployee() ([]internal.InboundOrdersPerEmployee, error) {
	args := m.Called()
	return args.Get(0).([]internal.InboundOrdersPerEmployee), args.Error(1)
}

func (m *MockEmployeeService) ReportInboundOrdersById(employeeId int) (internal.InboundOrdersPerEmployee, error) {
	args := m.Called(employeeId)
	return args.Get(0).(internal.InboundOrdersPerEmployee), args.Error(1)
}

type EmployeeTestSuiteMap struct {
	suite.Suite
	handler *handler.EmployeeHandlerDefault
	service *MockEmployeeService
}

func (e *EmployeeTestSuite) SetupTest() {
	var err error
	e.db, err = sql.Open("txdb_employee", "identier")
	require.NoError(e.T(), err)

	tempFile, err := os.CreateTemp("", "warehouse.json")
	require.NoError(e.T(), err)
	defer os.Remove(tempFile.Name()) // clean up

	warehouseJSON := `[{"warehouse_code": "WHS001", "address": "Rua Exemplo, 123", "telephone": "123456789", "minimum_capacity": 100, "minimum_temperature": 5.0}]`
	err = os.WriteFile(tempFile.Name(), []byte(warehouseJSON), 0644)
	require.NoError(e.T(), err)

	rp := repository.NewEmployeeMysql(e.db)
	rpWarehouse := repository.NewRepositoryWarehouse(nil, tempFile.Name())
	sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
	e.hd = handler.NewEmployeeDefault(sv)
}

func (e *EmployeeTestSuite) TestEmployeeHandler_GetAll() {
	defer e.db.Close()

	e.T().Run("should return all employees", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, employeeRoute, nil)
		w := httptest.NewRecorder()

		e.hd.GetAll(w, r)

		require.Equal(t, http.StatusOK, w.Result().StatusCode)

		var res struct {
			Data []internal.Employee `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&res)
		require.NoError(t, err)
		require.Len(t, res.Data, 10)
	})
}

func (e *EmployeeTestSuite) TestEmployeeHandler_Create() {
	defer e.db.Close()

	e.T().Run("should create a new employee", func(t *testing.T) {
		employee := internal.Employee{
			CardNumberId: "E1012",
			FirstName:    "Alice",
			LastName:     "Due",
			WarehouseId:  1,
		}

		body, err := json.Marshal(employee)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, employeeRoute, bytes.NewReader(body))
		w := httptest.NewRecorder()

		e.hd.Create(w, r)

		require.Equal(t, http.StatusCreated, w.Result().StatusCode)

		var res struct {
			Data internal.Employee `json:"data"`
		}
		err = json.NewDecoder(w.Body).Decode(&res)
		require.NoError(t, err)
		require.NotZero(t, res.Data.Id)
		require.Equal(t, "Alice", res.Data.FirstName)
	})

	e.T().Run("should return conflict if card_number_id already in use", func(t *testing.T) {

		employee := internal.Employee{
			CardNumberId: "E1001",
			FirstName:    "Bob",
			LastName:     "Marley",
			WarehouseId:  1,
		}

		body, err := json.Marshal(employee)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, employeeRoute, strings.NewReader(string(body)))
		w := httptest.NewRecorder()

		e.hd.Create(w, r)

		require.Equal(t, http.StatusConflict, w.Result().StatusCode)

		var res struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(w.Body).Decode(&res)
		require.NoError(t, err)
		require.Equal(t, "card number id already in use", res.Error)
	})
}

func (suite *EmployeeTestSuiteMap) TestGetEmployeeById() {
	employee := internal.Employee{Id: 1, CardNumberId: "E1001", FirstName: "John", LastName: "Doe", WarehouseId: 1}
	suite.service.On("GetById", 1).Return(employee, nil)

	r := httptest.NewRequest(http.MethodGet, "/employees/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.GetByID(w, r)

	assert.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)

	var response struct {
		Data internal.Employee `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), employee, response.Data)
}

func (suite *EmployeeTestSuiteMap) TestGetEmployeeByIdNotFound() {
	suite.service.On("GetById", 1).Return(internal.Employee{}, internal.ErrEmployeeNotFound)

	r := httptest.NewRequest(http.MethodGet, "/employees/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.GetByID(w, r)

	assert.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)

	var response struct {
		Error string `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), "employee not found", response.Error)
}

func (suite *EmployeeTestSuiteMap) TestUpdateEmployee() {
	employee := internal.Employee{Id: 1, CardNumberId: "E1001", FirstName: "John", LastName: "Doe", WarehouseId: 1}
	suite.service.On("Update", employee).Return(nil)
	suite.service.On("GetById", 1).Return(employee, nil)

	body, _ := json.Marshal(employee)
	r := httptest.NewRequest(http.MethodPut, "/employees/{id}", bytes.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.Update(w, r)

	assert.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)

	var response struct {
		Data internal.Employee `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), employee, response.Data)
}

func (suite *EmployeeTestSuiteMap) TestDeleteEmployee() {
	suite.service.On("Delete", 1).Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/employees/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.Delete(w, r)

	assert.Equal(suite.T(), http.StatusNoContent, w.Result().StatusCode)
}

func (e *EmployeeTestSuite) TestEmployeeHandler_ReportInboundOrders() {
	defer e.db.Close()

	e.T().Run("should return inbound orders count per employee", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, employeeRoute+"/report-inbound-orders", nil)
		w := httptest.NewRecorder()

		e.hd.ReportInboundOrders(w, r)

		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}

func TestEmployeeTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeTestSuite))
}
