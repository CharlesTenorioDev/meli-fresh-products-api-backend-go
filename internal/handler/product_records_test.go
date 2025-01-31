package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockProductRecordsService struct {
	mock.Mock
}

func (m *MockProductRecordsService) Create(product internal.ProductRecords) (internal.ProductRecords, error) {
	args := m.Called(product)
	product.ID = 1
	return args.Get(0).(internal.ProductRecords), args.Error(1)
}

func (m *MockProductRecordsService) GetAll() ([]internal.ProductRecords, error) {
	args := m.Called()
	return args.Get(0).([]internal.ProductRecords), args.Error(1)
}

func (m *MockProductRecordsService) GetByID(id int) (internal.ProductRecords, error) {
	args := m.Called(id)
	return args.Get(0).(internal.ProductRecords), args.Error(1)
}

func Test_ProductRecordsHandler_Create(t *testing.T) {
	type ResponseCreate struct {
		Data internal.ProductRecords `json:"data"`
	}
	tests := []struct {
		name           string
		mockSetup      func(*MockProductRecordsService)
		requestBody    interface{}
		expectedStatus int
		expectedBody   ResponseCreate
	}{
		{
			name: "Create_ok_status_200",
			mockSetup: func(p *MockProductRecordsService) {
				mockProduct := internal.ProductRecords{
					LastUpdateDate: time.Date(2025, 3, 21, 14, 30, 0, 0, time.Local),
					PurchasePrice:  100.0, SalePrice: 150.0,
				}
				mockCreatedProduct := internal.ProductRecords{
					LastUpdateDate: time.Date(2025, 3, 21, 14, 30, 0, 0, time.Local),
					PurchasePrice:  100.0, SalePrice: 150.0,
				}
				p.On("Create", mockProduct).Return(mockCreatedProduct, nil)
			},
			requestBody: internal.ProductRecords{
				LastUpdateDate: time.Date(2025, 3, 21, 14, 30, 0, 0, time.Local),
				PurchasePrice:  100.0, SalePrice: 150.0,
			},
			expectedStatus: http.StatusCreated,
			expectedBody: ResponseCreate{
				Data: internal.ProductRecords{
					LastUpdateDate: time.Date(2025, 3, 21, 14, 30, 0, 0, time.Local),
					PurchasePrice:  100.0, SalePrice: 150.0,
				},
			},
		},
		{
			name: "Invalid JSON, should return 422",
			mockSetup: func(p *MockProductRecordsService) {
			},
			requestBody:    "invalid-json",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   ResponseCreate{},
		},
		{
			name: "Error: Product unprocessable entity",
			mockSetup: func(p *MockProductRecordsService) {
				mockProduct := internal.ProductRecords{
					ProductID: 2, LastUpdateDate: time.Now().Truncate(24 * time.Hour),
					PurchasePrice: 200.0, SalePrice: 300.0,
				}
				p.On("Create", mockProduct).Return(internal.ProductRecords{}, internal.ErrProductUnprocessableEntity)
			},
			requestBody: internal.ProductRecords{
				ProductID: 2, LastUpdateDate: time.Now().Truncate(24 * time.Hour),
				PurchasePrice: 200.0, SalePrice: 300.0,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   ResponseCreate{},
		},
		{
			name: "Error: Product ID not found",
			mockSetup: func(p *MockProductRecordsService) {
				mockProduct := internal.ProductRecords{
					ProductID: 3, LastUpdateDate: time.Now().Truncate(24 * time.Hour),
					PurchasePrice: 150.0, SalePrice: 250.0,
				}
				p.On("Create", mockProduct).Return(internal.ProductRecords{}, internal.ErrProductIdNotFound)
			},
			requestBody: internal.ProductRecords{
				ProductID: 3, LastUpdateDate: time.Now().Truncate(24 * time.Hour),
				PurchasePrice: 150.0, SalePrice: 250.0,
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   ResponseCreate{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductRecordsService)
			productHandler := handler.NewProductRecordsDefault(mockService)
			tt.mockSetup(mockService)

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/productRecords", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			productHandler.Create(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			var actualResponse ResponseCreate

			err = json.NewDecoder(rec.Body).Decode(&actualResponse)
			require.NoError(t, err)
			fmt.Print(actualResponse)
			assert.Equal(t, actualResponse, tt.expectedBody)
		})
	}
}
