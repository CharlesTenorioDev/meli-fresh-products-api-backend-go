package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/require"
)

var jsonHeader = http.Header{
	"Content-Type": []string{"application/json"},
}

const (
	endpoint_warehouse  = "/api/v1/warehouses"
	file_path_warehouse = "warehouse.json"
)

type setup struct {
	router  *chi.Mux
	handler *WarehouseDefault
}

func newSetup(db map[int]repository.WarehouseAttributesMap) *setup {
	// Dependencies
	warehouseRepository := repository.NewRepositoryWarehouse(db, "")
	warehouseService := service.NewWarehouseDefault(warehouseRepository)
	warehouseHandler := NewWarehouseDefault(warehouseService)

	// Router
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	return &setup{
		router:  router,
		handler: warehouseHandler,
	}
}

func TestWarehouseHandler(t *testing.T) {
	db := map[int]repository.WarehouseAttributesMap{
		1: {WarehouseCode: "W1", Address: "123 Main St", Telephone: "123-456-7890", MinimumCapacity: 100, MinimumTemperature: -10.0},
		2: {WarehouseCode: "W2", Address: "456 Elm St", Telephone: "987-654-3210", MinimumCapacity: 200, MinimumTemperature: 0.0},
	}

	// Creating setup with a mock database
	setup := newSetup(db)

	// Router config
	router := chi.NewRouter()
	router.Route(endpoint_warehouse, func(r chi.Router) {
		r.Get("/", setup.handler.GetAll())
		r.Get("/{id}", setup.handler.GetByID())
		r.Post("/", setup.handler.Create())
		r.Patch("/{id}", setup.handler.Update())
		r.Delete("/{id}", setup.handler.Delete())
	})

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		expectedCode   int
		expectedHeader http.Header
		expectedBody   string
	}{
		// Successful tests cases
		{
			name:           "GetAllWarehouses",
			method:         "GET",
			url:            endpoint_warehouse,
			expectedCode:   http.StatusOK,
			expectedHeader: jsonHeader,
			expectedBody: `{"data":[
				{"id":1,"warehouse_code":"W1","address":"123 Main St","telephone":"123-456-7890","minimum_capacity":100,"minimum_temperature":-10},
				{"id":2,"warehouse_code":"W2","address":"456 Elm St","telephone":"987-654-3210","minimum_capacity":200,"minimum_temperature":0}
			]}`,
		},
		{
			name:           "GetWarehouseById",
			method:         "GET",
			url:            endpoint_warehouse + "/1",
			expectedCode:   http.StatusOK,
			expectedHeader: jsonHeader,
			expectedBody: `{"data":{
				"id":1,"warehouse_code":"W1","address":"123 Main St","telephone":"123-456-7890","minimum_capacity":100,"minimum_temperature":-10}}
			`,
		},
		{
			name:           "CreateWarehouse",
			method:         "POST",
			url:            endpoint_warehouse,
			body:           `{"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300,"minimum_temperature":5}`,
			expectedCode:   http.StatusCreated,
			expectedHeader: jsonHeader,
			expectedBody:   `{"data":{"id":3,"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300,"minimum_temperature":5}}`,
		},
		{
			name:           "UpdateWarehouse",
			method:         "PATCH",
			url:            endpoint_warehouse + "/1",
			body:           `{"address":"123 Main St UPDATED","telephone":"123-456-7890","minimum_capacity":1000,"minimum_temperature":-20}`,
			expectedCode:   http.StatusOK,
			expectedHeader: jsonHeader,
			expectedBody:   `{"data":{"id":1,"warehouse_code":"W1","address":"123 Main St UPDATED","telephone":"123-456-7890","minimum_capacity":1000,"minimum_temperature":-20}}`,
		},
		{
			name:           "DeleteWarehouse",
			method:         "DELETE",
			url:            endpoint_warehouse + "/1",
			expectedCode:   http.StatusNoContent,
			expectedHeader: http.Header{},
		},

		// Error test cases
		{
			name:           "GetWarehouseByIdInvalidID_Error",
			method:         "GET",
			url:            endpoint_warehouse + "/invalid",
			expectedCode:   http.StatusBadRequest,
			expectedHeader: jsonHeader,
			expectedBody:   `{"message":"Invalid ID format", "status":"Bad Request"}`,
		},
		{
			name:           "CreateWarehouseWithoutMinimumTemperature_Error",
			method:         "POST",
			url:            endpoint_warehouse,
			body:           `{"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300}`,
			expectedCode:   http.StatusUnprocessableEntity,
			expectedHeader: jsonHeader,
			expectedBody:   `{"message":"minimum temperature is required", "status":"Unprocessable Entity"}`,
		},
		{
			name:           "CreateWarehouseInvalidData_Error",
			method:         "POST",
			url:            endpoint_warehouse,
			body:           `{"warehouse_code":"W3","address":"789 Oak St","telephone":"555-1234","minimum_capacity":300,"minimum_temperature":"invalid"}`,
			expectedCode:   http.StatusBadRequest,
			expectedHeader: jsonHeader,
			expectedBody:   `{"message":"Invalid data", "status":"Bad Request"}`,
		},
		{
			name:           "UpdateWarehouseInvalidData_Error",
			method:         "PATCH",
			url:            endpoint_warehouse + "/1",
			body:           `{"address":"123 Main St UPDATED","telephone":"123-456-7890","minimum_capacity":1000,"minimum_temperature":"invalid"}`,
			expectedCode:   http.StatusBadRequest,
			expectedHeader: jsonHeader,
			expectedBody:   `{"message":"Invalid data", "status":"Bad Request"}`,
		},
		{
			name:           "UpdateWarehouseNotFound_Error",
			method:         "PATCH",
			url:            endpoint_warehouse + "/100",
			body:           `{"address":"123 Main St UPDATED","telephone":"123-456-7890","minimum_capacity":1000,"minimum_temperature":-20}`,
			expectedCode:   http.StatusNotFound,
			expectedHeader: jsonHeader,
			expectedBody:   `{"message":"Warehouse not found", "status":"Not Found"}`,
		},
		{
			name:           "DeleteWarehouseNotFound_Error",
			method:         "DELETE",
			url:            endpoint_warehouse + "/100",
			expectedCode:   http.StatusNotFound,
			expectedHeader: http.Header{},
		},
	}

	// Iterating over the tests
	for _, tt := range tests {
		fmt.Printf("Test: %s\n", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			// Request
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, tt.url, nil)
			}

			// Response
			res := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(res, req)

			// Validate the response
			if tt.method == "DELETE" {
				require.Equal(t, tt.expectedCode, res.Code)
				// Check if the warehouse was deleted
				_, ok := db[1]
				require.False(t, ok)
				return
			}

			require.Equal(t, tt.expectedCode, res.Code)
			require.Equal(t, tt.expectedHeader, res.Header())
			require.JSONEq(t, tt.expectedBody, res.Body.String())
		})
	}

}
