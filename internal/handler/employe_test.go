package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) GetAll() ([]internal.Employee, error) {
	args := m.Called()
	return args.Get(0).([]internal.Employee), nil
}

func (m *MockEmployeeService) GetById(id int) (internal.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *MockEmployeeService) Save(employee *internal.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockEmployeeService) Update(employee internal.Employee) error {
	args := m.Called(employee)
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

func (m *MockEmployeeService) ReportInboundOrdersById(employeeId int) (internal.InboundOrdersPerEmployee, error) {
	args := m.Called(employeeId)
	return args.Get(0).(internal.InboundOrdersPerEmployee), args.Error(1)
}

type EmployeeTestSuite struct {
	suite.Suite
	handler *EmployeeHandlerDefault
	service *MockEmployeeService
}

func (suite *EmployeeTestSuite) SetupTest() {
	suite.service = new(MockEmployeeService)
	suite.handler = NewEmployeeDefault(suite.service)
}

func (suite *EmployeeTestSuite) TestGetAllEmployees() {
	employees := map[int]internal.Employee{
		1: {Id: 1, CardNumberId: "E1001", FirstName: "John", LastName: "Doe", WarehouseId: 1},
	}
	suite.service.On("GetAll").Return(employees)

	r := httptest.NewRequest(http.MethodGet, "/employees", nil)
	w := httptest.NewRecorder()
	suite.handler.GetAll(w, r)

	assert.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)

	var response struct {
		Data map[int]internal.Employee `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)

	expected := map[int]internal.Employee{
		1: {Id: 1, CardNumberId: "E1001", FirstName: "John", LastName: "Doe", WarehouseId: 1},
	}
	assert.Equal(suite.T(), expected, response.Data)
}

func (suite *EmployeeTestSuite) TestGetEmployeeById() {
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

func (suite *EmployeeTestSuite) TestGetEmployeeByIdNotFound() {
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

func (suite *EmployeeTestSuite) TestSaveEmployee() {
	employee := internal.Employee{Id: 1, CardNumberId: "E1001", FirstName: "John", LastName: "Doe", WarehouseId: 1}
	suite.service.On("Save", &employee).Return(nil)

	body, _ := json.Marshal(employee)
	r := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewReader(body))
	w := httptest.NewRecorder()

	suite.handler.Create(w, r)

	assert.Equal(suite.T(), http.StatusCreated, w.Result().StatusCode)

	var response struct {
		Data internal.Employee `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), employee, response.Data)
}

func (suite *EmployeeTestSuite) TestSaveEmployeeError() {

	employee := internal.Employee{Id: 1, CardNumberId: "E1001", FirstName: "John", LastName: "Doe", WarehouseId: 1}

	suite.service.On("Save", &employee).Return(service.EmployeeInUse) //configure a mock

	// create the request body
	body, _ := json.Marshal(employee)
	r := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewReader(body))
	w := httptest.NewRecorder()

	suite.handler.Create(w, r) //calling the handler

	assert.Equal(suite.T(), http.StatusConflict, w.Result().StatusCode)

	var response struct {
		Error string `json:"error"` // The "error" field in the response
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response.Error)
	assert.Equal(suite.T(), "employee already in use", response.Error) // verifying the correct error message
}

func (suite *EmployeeTestSuite) TestUpdateEmployee() {
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

func (suite *EmployeeTestSuite) TestDeleteEmployee() {
	suite.service.On("Delete", 1).Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/employees/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.Delete(w, r)

	assert.Equal(suite.T(), http.StatusNoContent, w.Result().StatusCode)
}

func (suite *EmployeeTestSuite) TestDeleteEmployeeNotFound() {
	suite.service.On("Delete", 1).Return(internal.ErrEmployeeNotFound)

	r := httptest.NewRequest(http.MethodDelete, "/employees/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.Delete(w, r)

	// Verify the status code.
	assert.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)

	// Verify the message in the response body.
	var response struct {
		Error string `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "employee not found", response.Error)
}

func TestEmployeeTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeTestSuite))
}
