package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductBatchService struct {
	mock.Mock
}

type RequestProductBatchErrorJSON struct {
	BatchNumber        int     `json:"batch_number"`
	CurrentQuantity    int     `json:"current_quantity"`
	DueDate            string  `json:"due_date"`
	InitialQuantity    int     `json:"initial_quantity"`
	ManufacturingDate  string  `json:"manufacturing_date"`
	ManufacturingHour  int     `json:"manufacturing_hour"`
	MinumumTemperature float64 `json:"minumum_temperature"`
	ProductID          int     `json:"product_id"`
	SectionID          int     `json:"section_id"`
}

func (m *MockProductBatchService) FindByID(id int) (internal.ProductBatch, error) {
	args := m.Called(id)
	return args.Get(0).(internal.ProductBatch), args.Error(1)
}

func (m *MockProductBatchService) Save(prodBatch *internal.ProductBatch) error {
	if ok := prodBatch.Ok(); !ok {
		return internal.ErrProductBatchUnprocessableEntity
	}

	args := m.Called(prodBatch)
	prodBatch.ID = 1
	return args.Error(0)
}

func TestHandler_CreateProductBatchUnitTest(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockProductBatchService)
		requestBody        interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "successfully create a new prodBatch",
			mockSetup: func(m *MockProductBatchService) {
				m.On("Save", mock.Anything).Return(nil)
			},
			requestBody: handler.RequestProductBatchJSON{
				BatchNumber:        1234,
				CurrentQuantity:    100,
				CurrentTemperature: 40.5,
				DueDate:            "2022-01-08",
				InitialQuantity:    120,
				ManufacturingDate:  "2022-01-01 ",
				ManufacturingHour:  15,
				MinumumTemperature: -8,
				ProductID:          1,
				SectionID:          3,
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: internal.ProductBatch{
				ID:                 1,
				BatchNumber:        1234,
				CurrentQuantity:    100,
				CurrentTemperature: 40.5,
				DueDate:            "2022-01-08",
				InitialQuantity:    120,
				ManufacturingDate:  "2022-01-01 ",
				ManufacturingHour:  15,
				MinumumTemperature: -8,
				ProductID:          1,
				SectionID:          3,
			},
		},
		{
			name: "mandatory field error",
			mockSetup: func(m *MockProductBatchService) {
				m.On("Save", mock.Anything).Return(nil)
			},
			requestBody: RequestProductBatchErrorJSON{
				BatchNumber:        1234,
				CurrentQuantity:    100,
				DueDate:            "2022-01-08",
				InitialQuantity:    120,
				ManufacturingDate:  "2022-01-01 ",
				ManufacturingHour:  15,
				MinumumTemperature: -8,
				ProductID:          1,
				SectionID:          3,
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   *resterr.NewUnprocessableEntityError("current_temperature is required"),
		},
		{
			name:               "error while doing decode",
			mockSetup:          func(m *MockProductBatchService) {},
			requestBody:        "invalid json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   *resterr.NewBadRequestError("json: cannot unmarshal string into Go value of type map[string]interface {}"),
		},
		{
			name: "return fail error when required field is missing",
			mockSetup: func(m *MockProductBatchService) {
				m.On("Save", mock.Anything).Return(nil)
			},
			requestBody: handler.RequestProductBatchJSON{
				BatchNumber:        1234,
				CurrentQuantity:    100,
				CurrentTemperature: 40.5,
				DueDate:            "2022-01-08",
				InitialQuantity:    120,
				ManufacturingDate:  "2022-01-01 ",
				ManufacturingHour:  15,
				MinumumTemperature: -8,
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   *resterr.NewUnprocessableEntityError("couldn't parse product-batch"),
		},
		{
			name: "return conflict error when number is already in use",
			mockSetup: func(m *MockProductBatchService) {
				m.On("FindAll").Return([]internal.ProductBatch{
					{BatchNumber: 123},
				}, nil)
				m.On("Save", mock.Anything).Return(internal.ErrProductBatchNumberAlreadyInUse)
			},
			requestBody: handler.RequestProductBatchJSON{
				BatchNumber:        1234,
				CurrentQuantity:    100,
				CurrentTemperature: 40.5,
				DueDate:            "2022-01-08",
				InitialQuantity:    120,
				ManufacturingDate:  "2022-01-01 ",
				ManufacturingHour:  15,
				MinumumTemperature: -8,
				ProductID:          1,
				SectionID:          3,
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   *resterr.NewConflictError("product-batch with given product-batch number already registered"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductBatchService)
			sectionHandler := handler.NewHandlerProductBatch(mockService)
			tt.mockSetup(mockService)

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/product-batches", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("content-type", "application/json")

			rr := httptest.NewRecorder()
			hd := sectionHandler.Create
			hd(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case internal.ProductBatch:
					var actualResponse = struct {
						Data internal.ProductBatch `json:"data"`
					}{}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response, actualResponse.Data)
				case resterr.RestErr:
					var actualResponse resterr.RestErr
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response, actualResponse)
				default:
					t.Fatalf("Tipo de resposta inesperado: %T", response)
				}
			}
		})
	}
}

func TestHandler_ReadByIdProductBatchUnitTest(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockProductBatchService)
		id                 string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should return a product batch by ID",
			mockSetup: func(m *MockProductBatchService) {
				prodBatch := internal.ProductBatch{
					ID:                 1,
					BatchNumber:        1234,
					CurrentQuantity:    100,
					CurrentTemperature: 40.5,
					DueDate:            "2022-01-08",
					InitialQuantity:    120,
					ManufacturingDate:  "2022-01-01 ",
					ManufacturingHour:  15,
					MinumumTemperature: -8,
					ProductID:          1,
					SectionID:          3,
				}
				m.On("FindByID", 1).Return(prodBatch, nil)
			},
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": internal.ProductBatch{
					ID:                 1,
					BatchNumber:        1234,
					CurrentQuantity:    100,
					CurrentTemperature: 40.5,
					DueDate:            "2022-01-08",
					InitialQuantity:    120,
					ManufacturingDate:  "2022-01-01 ",
					ManufacturingHour:  15,
					MinumumTemperature: -8,
					ProductID:          1,
					SectionID:          3,
				},
			},
		},
		{
			name: "should return not found error",
			mockSetup: func(m *MockProductBatchService) {
				m.On("FindByID", 1).Return(internal.ProductBatch{}, internal.ErrProductBatchNotFound)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *resterr.NewNotFoundError("product-batch not found"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockProductBatchService) {
				m.On("FindByID", 1).Return(internal.ProductBatch{}, errors.New("internal server error"))
			},
			id:                 "1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name:               "should return bad request error",
			mockSetup:          func(m *MockProductBatchService) {},
			id:                 "invalid_id",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductBatchService)
			sectionHandler := handler.NewHandlerProductBatch(mockService)
			tt.mockSetup(mockService)

			req, err := http.NewRequest(http.MethodGet, "/product-batches/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			hd := sectionHandler.GetByID
			hd(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data internal.ProductBatch `json:"data"`
					}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)
				case resterr.RestErr:
					var actualResponse resterr.RestErr
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response, actualResponse)
				default:
					t.Fatalf("Tipo de resposta inesperado: %T", response)
				}
			}
		})
	}
}
