package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetAll() ([]internal.Product, error) {
	args := m.Called()
	return args.Get(0).([]internal.Product), args.Error(1)
}

func (m *MockProductService) GetByID(id int) (internal.Product, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (m *MockProductService) Create(product internal.Product) (internal.Product, error) {
	args := m.Called(product)
	product.Id = 1
	return args.Get(0).(internal.Product), args.Error(1)
}

func (m *MockProductService) Update(product internal.Product) (internal.Product, error) {
	args := m.Called(product)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (m *MockProductService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func Test_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedBody   map[string]any
	}{
		{
			name: "Success, show all products",
			mockSetup: func(p *MockProductService) {
				mockProduct := []internal.Product{
					{
						Id:                             1,
						ProductCode:                    "Product A",
						Description:                    "Test description",
						Height:                         10.0,
						Width:                          10.0,
						NetWeight:                      100,
						ExpirationRate:                 time.Now().Truncate(24 * time.Hour),
						RecommendedFreezingTemperature: 18,
						FreezingRate:                   18,
						ProductTypeId:                  1,
						SellerId:                       1,
					},
					{
						Id:                             2,
						ProductCode:                    "Product B",
						Description:                    "Test description",
						Height:                         10.0,
						Width:                          10.0,
						NetWeight:                      100,
						ExpirationRate:                 time.Now().Truncate(24 * time.Hour),
						RecommendedFreezingTemperature: 18,
						FreezingRate:                   18,
						ProductTypeId:                  1,
						SellerId:                       1,
					},
				}

				p.On("GetAll").Return(mockProduct, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]any{
				"data": []map[string]interface{}{
					{
						"id":                               1.0,
						"product_code":                     "Product A",
						"description":                      "Test description",
						"height":                           10.0,
						"width":                            10.0,
						"net_weight":                       100.0,
						"expiration_rate":                  time.Now().Truncate(24 * time.Hour).Format(time.RFC3339),
						"recommended_freezing_temperature": 18.0,
						"freezing_rate":                    18.0,
						"product_type_id":                  1.0,
						"seller_id":                        1.0,
					},
					{
						"id":                               2.0,
						"product_code":                     "Product B",
						"description":                      "Test description",
						"height":                           10.0,
						"width":                            10.0,
						"net_weight":                       100.0,
						"expiration_rate":                  time.Now().Truncate(24 * time.Hour).Format(time.RFC3339),
						"recommended_freezing_temperature": 18.0,
						"freezing_rate":                    18.0,
						"product_type_id":                  1.0,
						"seller_id":                        1.0,
					},
				},
			},
		},
		{
			name: "Error, return bad request 400.",
			mockSetup: func(p *MockProductService) {
				p.On("GetAll").Return([]internal.Product{}, errors.New("some error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			productHandler := handler.NewProducHandlerDefault(mockService)
			tt.mockSetup(mockService)

			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			rec := httptest.NewRecorder()

			productHandler.GetAll(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusOK {
				var body map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)

				for i, product := range body["data"].([]interface{}) {
					productMap := product.(map[string]interface{})
					expirationRate, ok := productMap["expiration_rate"].(string)
					if ok {
						_, err := time.Parse(time.RFC3339, expirationRate)
						assert.NoError(t, err)
					}
					assert.Equal(t, tt.expectedBody["data"].([]map[string]interface{})[i], productMap)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func Test_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockProductService)
		id             string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Success, get product by ID",
			mockSetup: func(p *MockProductService) {
				mockProduct := internal.Product{
					Id:                             1,
					ProductCode:                    "Product A",
					Description:                    "Test description",
					Height:                         10.0,
					Width:                          10.0,
					NetWeight:                      100,
					ExpirationRate:                 time.Now().Truncate(24 * time.Hour),
					RecommendedFreezingTemperature: 18,
					FreezingRate:                   18,
					ProductTypeId:                  1,
					SellerId:                       1,
				}
				p.On("GetByID", 1).Return(mockProduct, nil)
			},
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":                               1.0,
					"product_code":                     "Product A",
					"description":                      "Test description",
					"height":                           10.0,
					"width":                            10.0,
					"net_weight":                       100.0,
					"expiration_rate":                  time.Now().Truncate(24 * time.Hour).Format(time.RFC3339),
					"recommended_freezing_temperature": 18.0,
					"freezing_rate":                    18.0,
					"product_type_id":                  1.0,
					"seller_id":                        1.0,
				},
			},
		},
		{
			name: "Error, product not found",
			mockSetup: func(p *MockProductService) {
				p.On("GetByID", 1).Return(internal.Product{}, errors.New("product not found"))
			},
			id:             "1",
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"message": "product not found",
			},
		},
		{
			name: "Error, invalid ID format",
			mockSetup: func(p *MockProductService) {
				// No mock call needed for this case
			},
			id:             "invalid_id",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "strconv.Atoi: parsing \"invalid_id\": invalid syntax",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)                         // Cria mock do serviço
			productHandler := handler.NewProducHandlerDefault(mockService) // Cria handler
			tt.mockSetup(mockService)                                      // Configura comportamento do mock

			// Cria request com o ID como parâmetro
			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.id, nil)
			rec := httptest.NewRecorder()

			// Chama a função GetByID do handler
			productHandler.GetByID(rec, req)

			// Verifica o status da resposta
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// Verifica a estrutura da resposta JSON
			var body map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &body)
			assert.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody["data"], body["data"])
			} else {
				assert.Equal(t, tt.expectedBody["message"], body["message"])
			}

			// Garante que as expectativas do mock foram cumpridas
			mockService.AssertExpectations(t)
		})
	}
}
