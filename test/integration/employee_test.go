//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

func NewMockEmployeeService() *MockEmployeeService {
	return &MockEmployeeService{}
}

type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) GetAll() (db []internal.Employee, err error) {
	args := m.Called()
	return args.Get(0).([]internal.Employee), args.Error(1)
}

func (m *MockEmployeeService) GetByID(id int) (emp internal.Employee, err error) {
	args := m.Called(id)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *MockEmployeeService) Update(employees internal.Employee) (err error) {
	args := m.Called(employees)
	return args.Error(0)
}

func (m *MockEmployeeService) Save(emp *internal.Employee) (err error) {
	args := m.Called(emp)
	return args.Error(0)
}

func (m *MockEmployeeService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeService) CountInboundOrdersPerEmployee() ([]internal.InboundOrdersPerEmployee, error) {
	args := m.Called()
	return args.Get(0).([]internal.InboundOrdersPerEmployee), args.Error(1)
}

func (m *MockEmployeeService) ReportInboundOrdersByID(employeeID int) (io internal.InboundOrdersPerEmployee, err error) {
	args := m.Called(employeeID)
	return args.Get(0).(internal.InboundOrdersPerEmployee), args.Error(1)
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
			CardNumberID: "E1012",
			FirstName:    "Alice",
			LastName:     "Due",
			WarehouseID:  1,
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
		require.NotZero(t, res.Data.ID)
		require.Equal(t, "Alice", res.Data.FirstName)
	})

	e.T().Run("should return conflict if card_number_id already in use", func(t *testing.T) {

		employee := internal.Employee{
			CardNumberID: "E1001",
			FirstName:    "Bob",
			LastName:     "Marley",
			WarehouseID:  1,
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
	employee := internal.Employee{ID: 1, CardNumberID: "E1001", FirstName: "John", LastName: "Doe", WarehouseID: 1}
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
	employee := internal.Employee{ID: 1, CardNumberID: "E1001", FirstName: "John", LastName: "Doe", WarehouseID: 1}
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

func TestHandler_CreateEmployeeUnitTest(t *testing.T) {
	emp := internal.Employee{
		CardNumberID: "abcde",
		FirstName:    "Fabio",
		LastName:     "Nacarelli",
		WarehouseID:  14,
	}
	t.Run("create succeeds with 201", func(t *testing.T) {
		type CreatedRes struct {
			Data internal.Employee `json:"data"`
		}
		expectedStatus := http.StatusCreated
		expectedRes := CreatedRes{
			Data: emp,
		}
		sv := NewMockEmployeeService()
		sv.On("Save", &emp).Return(nil)
		b, _ := json.Marshal(emp)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			bytes.NewReader(b),
		)
		res := httptest.NewRecorder()
		hd.Create(res, req)

		var actualRes CreatedRes

		err := json.Unmarshal(res.Body.Bytes(), &actualRes)
		require.NoError(t, err)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
		sv.AssertNumberOfCalls(t, "Save", 1)
		require.Equal(t, expectedRes, actualRes)
	})
	t.Run("create fails with 422", func(t *testing.T) {
		expectedStatus := http.StatusUnprocessableEntity
		sv := NewMockEmployeeService()
		sv.On("Save", &emp).Return(errors.New("unprocessable entity"))
		hd := handler.NewEmployeeDefault(sv)
		b, _ := json.Marshal(emp)
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			bytes.NewReader(b),
		)
		res := httptest.NewRecorder()

		hd.Create(res, req)

		sv.AssertNumberOfCalls(t, "Save", 1)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("create fails with 409", func(t *testing.T) {
		expectedStatus := http.StatusConflict
		sv := NewMockEmployeeService()
		sv.On("Save", &emp).Return(service.ErrCardNumberIDInUse)
		hd := handler.NewEmployeeDefault(sv)
		b, _ := json.Marshal(emp)
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			bytes.NewReader(b),
		)
		res := httptest.NewRecorder()

		hd.Create(res, req)

		sv.AssertNumberOfCalls(t, "Save", 1)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("create with invalid body", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest
		sv := NewMockEmployeeService()
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			strings.NewReader(""),
		)
		res := httptest.NewRecorder()

		hd.Create(res, req)

		sv.AssertNumberOfCalls(t, "Save", 0)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("create fails because employee already in use", func(t *testing.T) {
		expectedStatus := http.StatusConflict
		sv := NewMockEmployeeService()
		sv.On("Save", &emp).Return(service.ErrEmployeeInUse)
		hd := handler.NewEmployeeDefault(sv)
		b, _ := json.Marshal(emp)
		req := httptest.NewRequest(
			http.MethodPost,
			"/",
			bytes.NewReader(b),
		)
		res := httptest.NewRecorder()

		hd.Create(res, req)

		sv.AssertNumberOfCalls(t, "Save", 1)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
}

func TestHandler_ReadEmployeeUnitTest(t *testing.T) {
	employeeDb := []internal.Employee{
		{
			ID:           1,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			CardNumberID: "CN001",
			WarehouseID:  12,
		},
		{
			ID:           2,
			FirstName:    "Mocked",
			LastName:     "Database",
			CardNumberID: "CN002",
			WarehouseID:  11,
		},
	}
	t.Run("fetch every employee", func(t *testing.T) {
		type GetAllRes struct {
			Data []internal.Employee `json:"data"`
		}

		expectedStatus := http.StatusOK
		expectedRes := GetAllRes{
			Data: employeeDb,
		}
		sv := NewMockEmployeeService()
		sv.On("GetAll").Return(employeeDb, nil)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		hd.GetAll(res, req)

		var actualRes GetAllRes
		err := json.Unmarshal(res.Body.Bytes(), &actualRes)
		require.NoError(t, err)
		sv.AssertNumberOfCalls(t, "GetAll", 1)
		require.Equal(t, expectedRes, actualRes)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("fetch employee by id (404)", func(t *testing.T) {
		type GetByIDRes struct {
			Data string `json:"data"`
		}

		expectedStatus := http.StatusNotFound
		expectedRes := GetByIDRes{
			Data: "employee not found",
		}
		sv := NewMockEmployeeService()
		sv.On("GetByID", 1).Return(internal.Employee{}, errors.New("employee not found"))
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.GetByID(res, req)

		var actualRes GetByIDRes
		err := json.Unmarshal(res.Body.Bytes(), &actualRes)
		require.NoError(t, err)
		sv.AssertNumberOfCalls(t, "GetByID", 1)
		require.Equal(t, expectedRes, actualRes)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("fetch employee by id (200)", func(t *testing.T) {
		employee := internal.Employee{
			ID:           1,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			CardNumberID: "FN001",
			WarehouseID:  14,
		}
		type GetByIDRes struct {
			Data internal.Employee `json:"data"`
		}

		expectedStatus := http.StatusOK
		expectedRes := GetByIDRes{
			Data: employee,
		}
		sv := NewMockEmployeeService()
		sv.On("GetByID", 1).Return(employee, nil)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.GetByID(res, req)

		var actualRes GetByIDRes
		err := json.Unmarshal(res.Body.Bytes(), &actualRes)
		require.NoError(t, err)
		sv.AssertNumberOfCalls(t, "GetByID", 1)
		require.Equal(t, expectedRes, actualRes)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("fetching every employee fails", func(t *testing.T) {
		expectedStatus := http.StatusInternalServerError
		sv := NewMockEmployeeService()
		sv.On("GetAll").Return([]internal.Employee{}, errors.New("internal server error"))
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		hd.GetAll(res, req)

		sv.AssertNumberOfCalls(t, "GetAll", 1)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("fetch employee by id (invalid)", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest
		sv := NewMockEmployeeService()
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "abcdef")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.GetByID(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
}

func TestHandler_UpdateEmployeeUnitTest(t *testing.T) {
	t.Run("update successfully (200)", func(t *testing.T) {
		employee := internal.Employee{
			ID:           1,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			CardNumberID: "FN001",
			WarehouseID:  14,
		}
		type UpdateRes struct {
			Data internal.Employee `json:"data"`
		}

		expectedStatus := http.StatusOK
		expectedRes := UpdateRes{
			Data: employee,
		}
		sv := NewMockEmployeeService()
		sv.On("Update", employee).Return(nil)
		sv.On("GetByID", 1).Return(employee, nil)
		b, _ := json.Marshal(employee)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodPatch, "/{id}", bytes.NewReader(b))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.Update(res, req)

		var actualRes UpdateRes
		err := json.Unmarshal(res.Body.Bytes(), &actualRes)
		require.NoError(t, err)
		sv.AssertNumberOfCalls(t, "Update", 1)
		sv.AssertNumberOfCalls(t, "GetByID", 1)
		require.Equal(t, expectedRes, actualRes)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("update fails (404)", func(t *testing.T) {
		employee := internal.Employee{
			ID:           1,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			CardNumberID: "FN001",
			WarehouseID:  14,
		}
		type UpdateRes struct {
			Data string `json:"data"`
		}

		expectedStatus := http.StatusNotFound
		expectedRes := UpdateRes{
			Data: service.ErrEmployeeNotFound.Error(),
		}
		sv := NewMockEmployeeService()
		sv.On("Update", employee).Return(service.ErrEmployeeNotFound)
		b, _ := json.Marshal(employee)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodPatch, "/{id}", bytes.NewReader(b))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.Update(res, req)

		var actualRes UpdateRes
		err := json.Unmarshal(res.Body.Bytes(), &actualRes)
		sv.AssertNumberOfCalls(t, "Update", 1)
		require.NoError(t, err)
		require.Equal(t, expectedRes, actualRes)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("update fails because invalid id", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest
		sv := NewMockEmployeeService()
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodPatch, "/{id}", strings.NewReader(""))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "abcdef")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.Update(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("update fails with conflict", func(t *testing.T) {
		employee := internal.Employee{
			ID:           1,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			CardNumberID: "FN001",
			WarehouseID:  14,
		}
		type UpdateRes struct {
			Data string `json:"data"`
		}

		expectedStatus := http.StatusConflict
		expectedRes := UpdateRes{
			Data: service.ErrConflictInEmployee.Error(),
		}
		sv := NewMockEmployeeService()
		sv.On("Update", employee).Return(service.ErrConflictInEmployee)
		b, _ := json.Marshal(employee)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodPatch, "/{id}", bytes.NewReader(b))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.Update(res, req)

		var actualRes UpdateRes
		err := json.Unmarshal(res.Body.Bytes(), &actualRes)
		require.NoError(t, err)
		sv.AssertNumberOfCalls(t, "Update", 1)
		require.Equal(t, expectedRes, actualRes)
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("update fails, bad body", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest
		sv := NewMockEmployeeService()
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodPatch, "/{id}", strings.NewReader(""))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.Update(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("update fails, internal get by id fails", func(t *testing.T) {
		employee := internal.Employee{
			ID:           1,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			CardNumberID: "FN001",
			WarehouseID:  14,
		}

		expectedStatus := http.StatusInternalServerError
		sv := NewMockEmployeeService()
		sv.On("Update", employee).Return(nil)
		sv.On("GetByID", 1).Return(internal.Employee{}, errors.New("error retrieving updated employee"))
		b, _ := json.Marshal(employee)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodPatch, "/{id}", bytes.NewReader(b))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.Update(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
}

func TestHandler_DeleteEmployeeUnitTest(t *testing.T) {
	t.Run("delete successfully (204)", func(t *testing.T) {
		expectedStatus := http.StatusNoContent
		sv := NewMockEmployeeService()
		sv.On("Delete", 1).Return(nil)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodDelete, "/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.Delete(res, req)

		sv.AssertNumberOfCalls(t, "Delete", 1)
		require.Equal(t, 0, res.Body.Len())
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("delete fails (404)", func(t *testing.T) {
		expectedStatus := http.StatusNotFound
		expectedRes := `{"data":"employee not found"}`
		sv := NewMockEmployeeService()
		sv.On("Delete", 1).Return(errors.New("employee not found"))
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodDelete, "/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.Delete(res, req)

		sv.AssertNumberOfCalls(t, "Delete", 1)
		require.Equal(t, expectedRes, string(res.Body.Bytes()))
		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("delete fails because of invalid id", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest
		sv := NewMockEmployeeService()
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodDelete, "/{id}", strings.NewReader(""))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "abcdef")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()

		hd.Delete(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
}

func TestHandler_ReportInboundOrdersEmployeeUnitTest(t *testing.T) {
	t.Run("fails to fetch inbound orders", func(t *testing.T) {
		expectedStatus := http.StatusInternalServerError
		sv := NewMockEmployeeService()
		sv.On("CountInboundOrdersPerEmployee").Return([]internal.InboundOrdersPerEmployee{}, errors.New("failed to fetch inbound orders"))
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		hd.ReportInboundOrders(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("succeeds to fetch inbound orders", func(t *testing.T) {
		expectedStatus := http.StatusOK
		sv := NewMockEmployeeService()
		sv.On("CountInboundOrdersPerEmployee").Return([]internal.InboundOrdersPerEmployee{}, nil)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		hd.ReportInboundOrders(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("fails with invalid id", func(t *testing.T) {
		expectedStatus := http.StatusBadRequest
		sv := NewMockEmployeeService()
		sv.On("CountInboundOrdersPerEmployee").Return([]internal.InboundOrdersPerEmployee{}, nil)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/?id=abcdef", nil)
		res := httptest.NewRecorder()

		hd.ReportInboundOrders(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("employee not found", func(t *testing.T) {
		expectedStatus := http.StatusNotFound
		sv := NewMockEmployeeService()
		sv.On("ReportInboundOrdersByID", 1).Return(internal.InboundOrdersPerEmployee{}, errors.New("employee not found"))
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/?id=1", nil)
		res := httptest.NewRecorder()

		hd.ReportInboundOrders(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
	t.Run("employee found", func(t *testing.T) {
		expectedStatus := http.StatusOK
		sv := NewMockEmployeeService()
		sv.On("ReportInboundOrdersByID", 1).Return(internal.InboundOrdersPerEmployee{}, nil)
		hd := handler.NewEmployeeDefault(sv)
		req := httptest.NewRequest(http.MethodGet, "/?id=1", nil)
		res := httptest.NewRecorder()

		hd.ReportInboundOrders(res, req)

		require.Equal(t, expectedStatus, res.Result().StatusCode)
	})
}
func TestEmployeeTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeTestSuite))
}
