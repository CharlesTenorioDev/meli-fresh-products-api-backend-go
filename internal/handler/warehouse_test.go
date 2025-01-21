package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func NewWarehouseServiceMock() *WarehouseServiceMock {
	return &WarehouseServiceMock{}
}

// Warehouse Service Mock
type WarehouseServiceMock struct {
	mock.Mock
}

// Warehouse Service FindAll returns a list of warehouses
func (w *WarehouseServiceMock) FindAll() ([]internal.Warehouse, error) {
	args := w.Called()
	return args.Get(0).([]internal.Warehouse), args.Error(1)
}

// Warehouse Service FindByID returns a warehouse by id
func (w *WarehouseServiceMock) FindByID(id int) (internal.Warehouse, error) {
	args := w.Called(id)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

// Warehouse Service Save creates a new warehouse
func (w *WarehouseServiceMock) Save(warehouse *internal.Warehouse) error {
	args := w.Called(warehouse)
	return args.Error(0)
}

// Warehouse Service Update updates a warehouse
func (w *WarehouseServiceMock) Update(id int, warehousePatch *internal.WarehousePatchUpdate) (internal.Warehouse, error) {
	args := w.Called(warehousePatch)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

// Warehouse Service Delete deletes a warehouse by id
func (w *WarehouseServiceMock) Delete(id int) error {
	args := w.Called(id)
	return args.Error(0)
}

var (
	endpointWarehouse = "/api/v1/warehouses"
	jsonHeader        = http.Header{"Content-Type": []string{"application/json"}}
)

// TestCases is a struct that contains the test cases for the warehouse handler
type TestCases struct {
	description       string
	method            string
	url               string
	id                string
	body              string
	expectedBody      string
	expectedCode      int
	expectedHeader    http.Header
	mock              func() *WarehouseServiceMock
	expectedMockCalls int
}

func TestWarehouseHandler_Create(t *testing.T) {
	cases := []*TestCases{
		{
			description:    "case 1 - success: Create a new warehouse",
			method:         "POST",
			url:            endpointWarehouse,
			body:           `{"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300,"minimum_temperature":5}`,
			expectedBody:   `{"data":{"id":3,"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300,"minimum_temperature":5}}`,
			expectedCode:   http.StatusCreated,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Save", mock.AnythingOfType("*internal.Warehouse")).Run(func(args mock.Arguments) {
					w := args.Get(0).(*internal.Warehouse)
					w.ID = 3
				}).Return(nil)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 2 - error: Attempt to create a new warehouse without minimum temperature",
			method:      "POST",
			url:         endpointWarehouse,
			body:        `{"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300}`,
			expectedBody: `{
				"message": "unprocessable entity: minimum temperature is required",
				"error": "unprocessable_entity",
				"code": 422,
				"causes": null
			}`,
			expectedCode:   http.StatusUnprocessableEntity,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Save", mock.Anything).Return(fmt.Errorf("%w: %v", internal.ErrWarehouseUnprocessableEntity, "minimum temperature is required"))
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 3 - error: Attempt to create a new warehouse without minimum capacity",
			method:      "POST",
			url:         endpointWarehouse,
			body:        `{"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_temperature":5}`,
			expectedBody: `{
				"message": "unprocessable entity: minimum capacity is required",
				"error": "unprocessable_entity",
				"code": 422,
				"causes": null
			}`,
			expectedCode:   http.StatusUnprocessableEntity,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Save", mock.Anything).Return(fmt.Errorf("%w: %v", internal.ErrWarehouseUnprocessableEntity, "minimum capacity is required"))
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 4 - error: Attempt to create a new warehouse with invalid data",
			method:      "POST",
			url:         endpointWarehouse,
			body:        `{"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300,"minimum_temperature":"invalid"}`,
			expectedBody: `{
				"message": "Invalid data",
				"error": "bad_request",
				"code": 400,
				"causes": null
			}`,
			expectedCode:   http.StatusBadRequest,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				return mk
			},
			expectedMockCalls: 0,
		},
		{
			description: "case 5 - error: Attempt to create a new warehouse with an existing warehouse code",
			method:      "POST",
			url:         endpointWarehouse,
			body:        `{"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300,"minimum_temperature":5}`,
			expectedBody: `{
				"message": "warehouse already exists",
				"error": "conflict",
				"code": 409,
				"causes": null
			}`,
			expectedCode:   http.StatusConflict,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Save", mock.Anything).Return(internal.ErrWarehouseRepositoryDuplicated)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 6 - error: Attempt to create a new warehouse generating an unexpected error",
			method:      "POST",
			url:         endpointWarehouse,
			body:        `{"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300,"minimum_temperature":5}`,
			expectedBody: `{
				"message": "Internal Server Error",
				"error": "internal_server_error",
				"code": 500,
				"causes": null
			}`,
			expectedCode:   http.StatusInternalServerError,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Save", mock.Anything).Return(errors.New("unexpected error"))
				return mk
			},
			expectedMockCalls: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			// handler and its dependencies
			sv := tc.mock()
			hd := handler.NewWarehouseDefault(sv)
			hdFunc := hd.Create()

			// http request and response
			request := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// WHEN
			hdFunc(response, request)

			// THEN
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			sv.AssertNumberOfCalls(t, "Save", tc.expectedMockCalls)
		})
	}
}

func TestWarehouseHandler_GetAll(t *testing.T) {
	cases := []*TestCases{
		{
			description:    "case 1 - success: Get all warehouses",
			method:         "GET",
			url:            endpointWarehouse,
			expectedCode:   http.StatusOK,
			expectedHeader: jsonHeader,
			expectedBody: `{"data":[
				{"id":1,"warehouse_code":"W1","address":"123 Main St","telephone":"123-456-7890","minimum_capacity":100,"minimum_temperature":-10},
				{"id":2,"warehouse_code":"W2","address":"456 Elm St","telephone":"987-654-3210","minimum_capacity":200,"minimum_temperature":0}
			]}`,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("FindAll").Return([]internal.Warehouse{
					{ID: 1, WarehouseCode: "W1", Address: "123 Main St", Telephone: "123-456-7890", MinimumCapacity: 100, MinimumTemperature: -10},
					{ID: 2, WarehouseCode: "W2", Address: "456 Elm St", Telephone: "987-654-3210", MinimumCapacity: 200, MinimumTemperature: 0},
				}, nil)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description:    "case 2 - error: Attempt to get all warehouses generating an unexpected error",
			method:         "GET",
			url:            endpointWarehouse,
			expectedCode:   http.StatusInternalServerError,
			expectedHeader: jsonHeader,
			expectedBody: `{
				"message": "Internal Server Error",
				"error": "internal_server_error",
				"code": 500,
				"causes": null
			}`,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("FindAll").Return([]internal.Warehouse{}, errors.New("unexpected error"))
				return mk
			},
			expectedMockCalls: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			// handler and its dependencies
			sv := tc.mock()
			hd := handler.NewWarehouseDefault(sv)
			hdFunc := hd.GetAll()

			// http request and response
			request := httptest.NewRequest(tc.method, tc.url, nil)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// WHEN
			hdFunc(response, request)

			// THEN
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			sv.AssertNumberOfCalls(t, "FindAll", tc.expectedMockCalls)
		})
	}
}

func TestWarehouseHandler_GetByID(t *testing.T) {
	cases := []*TestCases{
		{
			description:  "case 1 - success: Get a warehouse by id",
			method:       "GET",
			url:          endpointWarehouse,
			id:           "1",
			expectedCode: http.StatusOK,
			expectedBody: `{"data":{
				"id":1,"warehouse_code":"W1","address":"123 Main St","telephone":"123-456-7890","minimum_capacity":100,"minimum_temperature":-10}
			}`,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("FindByID", 1).Return(internal.Warehouse{
					ID: 1, WarehouseCode: "W1", Address: "123 Main St", Telephone: "123-456-7890", MinimumCapacity: 100, MinimumTemperature: -10,
				}, nil)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description:  "case 2 - error: Attempt to get a warehouse by a non existent id",
			method:       "GET",
			url:          endpointWarehouse,
			id:           "100",
			expectedCode: http.StatusNotFound,
			expectedBody: `{
				"message": "warehouse not found",
				"error": "not_found",
				"code": 404,
				"causes": null
			}`,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("FindByID", 100).Return(internal.Warehouse{}, internal.ErrWarehouseRepositoryNotFound)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description:  "case 3 - error: Attempt to get a warehouse by an invalid id",
			method:       "GET",
			url:          endpointWarehouse,
			id:           "invalid",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{
				"message": "Invalid ID format",
				"error": "bad_request",
				"code": 400,
				"causes": null
			}`,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				return mk
			},
			expectedMockCalls: 0,
		},
		{
			description:  "case 4 - error: Attempt to get a warehouse generating an unexpected error",
			method:       "GET",
			url:          endpointWarehouse,
			id:           "1",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{
				"message": "Internal Server Error",
				"error": "internal_server_error",
				"code": 500,
				"causes": null
			}`,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("FindByID", 1).Return(internal.Warehouse{}, errors.New("unexpected error"))
				return mk
			},
			expectedMockCalls: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			// handler and its dependencies
			sv := tc.mock()
			hd := handler.NewWarehouseDefault(sv)
			hdFunc := hd.GetByID()

			// http request and response
			request := httptest.NewRequest(tc.method, tc.url+"/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// WHEN
			hdFunc(response, request)

			// THEN
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			sv.AssertNumberOfCalls(t, "FindByID", tc.expectedMockCalls)
		})
	}
}

func TestWarehouseHandler_Update(t *testing.T) {
	cases := []*TestCases{
		{
			description:    "case 1 - success: Update a warehouse",
			method:         "PATCH",
			url:            endpointWarehouse,
			id:             "1",
			body:           `{"address":"123 Main St UPDATED","telephone":"123-456-7890","minimum_capacity":1000,"minimum_temperature":-20}`,
			expectedBody:   `{"data":{"id":1,"warehouse_code":"W1","address":"123 Main St UPDATED","telephone":"123-456-7890","minimum_capacity":1000,"minimum_temperature":-20}}`,
			expectedCode:   http.StatusOK,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Update", mock.AnythingOfType("*internal.WarehousePatchUpdate")).Return(internal.Warehouse{
					ID: 1, WarehouseCode: "W1", Address: "123 Main St UPDATED", Telephone: "123-456-7890", MinimumCapacity: 1000, MinimumTemperature: -20,
				}, nil)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 2 - error: Attempt to update a non existent warehouse",
			method:      "PATCH",
			url:         endpointWarehouse,
			id:          "100",
			body:        `{"address":"123 Main St UPDATED","telephone":"123-456-7890","minimum_capacity":1000,"minimum_temperature":-20}`,
			expectedBody: `{
				"message": "warehouse not found",
				"error": "not_found",
				"code": 404,
				"causes": null
			}`,
			expectedCode:   http.StatusNotFound,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Update", mock.AnythingOfType("*internal.WarehousePatchUpdate")).Return(internal.Warehouse{}, internal.ErrWarehouseRepositoryNotFound)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 3 - error: Attempt to update a warehouse with invalid data",
			method:      "PATCH",
			url:         endpointWarehouse,
			id:          "1",
			body:        `{"address":"123 Main St UPDATED","telephone":"123-456-7890","minimum_capacity":1000,"minimum_temperature":"invalid"}`,
			expectedBody: `{
				"message": "Invalid data",
				"error": "bad_request",
				"code": 400,
				"causes": null
			}`,
			expectedCode:   http.StatusBadRequest,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				return mk
			},
			expectedMockCalls: 0,
		},
		{
			description: "case 4 - error: Attempt to update a warehouse generating an unexpected error",
			method:      "PATCH",
			url:         endpointWarehouse,
			id:          "1",
			body:        `{"address":"123 Main St UPDATED","telephone":"123-456-7890","minimum_capacity":1000,"minimum_temperature":-20}`,
			expectedBody: `{
				"message": "Internal Server Error",
				"error": "internal_server_error",
				"code": 500,
				"causes": null
			}`,
			expectedCode:   http.StatusInternalServerError,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Update", mock.AnythingOfType("*internal.WarehousePatchUpdate")).Return(internal.Warehouse{}, errors.New("unexpected error"))
				return mk
			},
			expectedMockCalls: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			// handler and its dependencies
			sv := tc.mock()
			hd := handler.NewWarehouseDefault(sv)
			hdFunc := hd.Update()

			// http request and response
			request := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// WHEN
			hdFunc(response, request)

			// THEN
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			sv.AssertNumberOfCalls(t, "Update", tc.expectedMockCalls)
		})
	}
}

func TestWarehouseHandler_Delete(t *testing.T) {
	cases := []*TestCases{
		{
			description:    "case 1 - success: Delete a warehouse",
			method:         "DELETE",
			url:            endpointWarehouse,
			id:             "1",
			expectedCode:   http.StatusNoContent,
			expectedHeader: http.Header{},
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Delete", 1).Return(nil)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description:  "case 2 - error: Attempt to delete a non existent warehouse",
			method:       "DELETE",
			url:          endpointWarehouse,
			id:           "100",
			expectedCode: http.StatusNotFound,
			expectedBody: `{
				"message": "warehouse not found",
				"error": "not_found",
				"code": 404,
				"causes": null
			}`,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Delete", 100).Return(internal.ErrWarehouseRepositoryNotFound)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description:  "case 3 - error: Attempt to delete a warehouse generating an unexpected error",
			method:       "DELETE",
			url:          endpointWarehouse,
			id:           "1",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{
				"message": "Internal Server Error",
				"error": "internal_server_error",
				"code": 500,
				"causes": null
			}`,
			expectedHeader: jsonHeader,
			mock: func() *WarehouseServiceMock {
				mk := NewWarehouseServiceMock()
				mk.On("Delete", 1).Return(errors.New("unexpected error"))
				return mk
			},
			expectedMockCalls: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			// handler and its dependencies
			sv := tc.mock()
			hd := handler.NewWarehouseDefault(sv)
			hdFunc := hd.Delete()

			// http request and response
			request := httptest.NewRequest(tc.method, tc.url, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// WHEN
			hdFunc(response, request)

			// THEN
			require.Equal(t, tc.expectedCode, response.Code)
			if tc.expectedBody != "" {
				require.JSONEq(t, tc.expectedBody, response.Body.String())
			}
			sv.AssertNumberOfCalls(t, "Delete", tc.expectedMockCalls)
		})
	}
}
