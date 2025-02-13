package handler_test

import (
	"bytes"
	"errors"
	"fmt"

	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	product.ID = 1
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

func (m *MockProductService) GetAllRecord() ([]internal.ProductRecordsJSONCount, error) {
	args := m.Called()
	return args.Get(0).([]internal.ProductRecordsJSONCount), args.Error(1)
}

func (m *MockProductService) GetByIDRecord(id int) (internal.ProductRecordsJSONCount, error) {
	args := m.Called(id)
	return args.Get(0).(internal.ProductRecordsJSONCount), args.Error(1)
}

func Test_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Find_All_status_200",
			mockSetup: func(p *MockProductService) {
				mockProduct := []internal.Product{
					{
						ID:                             1,
						ProductCode:                    "Product A",
						Description:                    "Test description",
						Height:                         10.0,
						Width:                          10.0,
						NetWeight:                      100,
						ExpirationRate:                 1,
						RecommendedFreezingTemperature: 18,
						FreezingRate:                   18,
						ProductTypeID:                  1,
						SellerID:                       1,
					},
					{
						ID:                             2,
						ProductCode:                    "Product B",
						Description:                    "Test description",
						Height:                         10.0,
						Width:                          10.0,
						NetWeight:                      100,
						ExpirationRate:                 1,
						RecommendedFreezingTemperature: 18,
						FreezingRate:                   18,
						ProductTypeID:                  1,
						SellerID:                       1,
					},
				}
				p.On("GetAll").Return(mockProduct, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []internal.Product{
					{
						ID:                             1,
						ProductCode:                    "Product A",
						Description:                    "Test description",
						Height:                         10.0,
						Width:                          10.0,
						NetWeight:                      100,
						ExpirationRate:                 1,
						RecommendedFreezingTemperature: 18,
						FreezingRate:                   18,
						ProductTypeID:                  1,
						SellerID:                       1,
					},
					{
						ID:                             2,
						ProductCode:                    "Product B",
						Description:                    "Test description",
						Height:                         10.0,
						Width:                          10.0,
						NetWeight:                      100,
						ExpirationRate:                 1,
						RecommendedFreezingTemperature: 18,
						FreezingRate:                   18,
						ProductTypeID:                  1,
						SellerID:                       1,
					},
				},
			},
		},
		{
			name: "Find_All_status_500",
			mockSetup: func(p *MockProductService) {
				p.On("GetAll").Return([]internal.Product{}, errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "Find_All_status_404",
			mockSetup: func(p *MockProductService) {
				p.On("GetAll").Return([]internal.Product{}, internal.ErrProductNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			productHandler := handler.NewProductHandlerDefault(mockService)
			tt.mockSetup(mockService)

			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			rec := httptest.NewRecorder()

			productHandler.GetAll(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != nil {
				switch response := tt.expectedBody.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data []internal.Product `json:"data"`
					}
					err := json.NewDecoder(rec.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)

				case resterr.RestErr:
					var actualResponse resterr.RestErr
					err := json.NewDecoder(rec.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response, actualResponse)
				default:
					t.Fatalf("Tipo de resposta inesperado: %T", response)
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
		expectedBody   interface{}
	}{
		{
			name: "find_by_id_existent_status_200",
			mockSetup: func(p *MockProductService) {
				mockProduct := internal.Product{
					ID:                             1,
					ProductCode:                    "Product A",
					Description:                    "Test description",
					Height:                         10.0,
					Width:                          10.0,
					NetWeight:                      100,
					ExpirationRate:                 1,
					RecommendedFreezingTemperature: 18,
					FreezingRate:                   18,
					ProductTypeID:                  1,
					SellerID:                       1,
				}
				p.On("GetByID", 1).Return(mockProduct, nil)
			},
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": internal.Product{
					ID:                             1,
					ProductCode:                    "Product A",
					Description:                    "Test description",
					Height:                         10.0,
					Width:                          10.0,
					NetWeight:                      100,
					ExpirationRate:                 1,
					RecommendedFreezingTemperature: 18,
					FreezingRate:                   18,
					ProductTypeID:                  1,
					SellerID:                       1,
				},
			},
		},
		{
			name: "find_by_id_non_existent_status_404",
			mockSetup: func(p *MockProductService) {
				p.On("GetByID", 1).Return(internal.Product{}, errors.New("product not found"))
			},
			id:             "1",
			expectedStatus: http.StatusNotFound,
			expectedBody:   *resterr.NewNotFoundError("product not found"),
		},
		{
			name:           "find_by_id_status_400",
			mockSetup:      func(p *MockProductService) {},
			id:             "invalid_id",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			productHandler := handler.NewProductHandlerDefault(mockService)
			tt.mockSetup(mockService)

			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.id, nil)
			rec := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			productHandler.GetByID(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != nil {
				switch response := tt.expectedBody.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data internal.Product `json:"data"`
					}
					err := json.NewDecoder(rec.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)
				case resterr.RestErr:
					var actualResponse resterr.RestErr
					err := json.NewDecoder(rec.Body).Decode(&actualResponse)
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

func Test_Create(t *testing.T) {
	type ResponseCreate struct {
		Data internal.Product `json:"data"`
	}
	tests := []struct {
		name             string
		mockSetup        func(*MockProductService)
		requestBody      interface{}
		expectedStatus   int
		expectedResponse ResponseCreate
	}{
		{
			name: "create_ok_status_201",
			mockSetup: func(m *MockProductService) {
				m.On("Create", mock.Anything).Return(internal.Product{
					ID:                             1,
					ProductCode:                    "Product A",
					Description:                    "Test description",
					Height:                         10.0,
					Length:                         10.0, // Adicionado aqui se importar
					NetWeight:                      100.0,
					ExpirationRate:                 1.0,
					RecommendedFreezingTemperature: 18.0,
					FreezingRate:                   18.0,
					ProductTypeID:                  1,
					SellerID:                       1,
					Width:                          10.0,
				}, nil)
			},
			requestBody: internal.Product{
				ProductCode:                    "Product A",
				Description:                    "Test description",
				Height:                         10.0,
				Length:                         10.0,
				NetWeight:                      100,
				ExpirationRate:                 1,
				RecommendedFreezingTemperature: 18,
				FreezingRate:                   18,
				ProductTypeID:                  1,
				SellerID:                       1,
				Width:                          10.0,
			},
			expectedStatus: http.StatusCreated,
			expectedResponse: ResponseCreate{
				Data: internal.Product{
					ProductCode:                    "Product A",
					Description:                    "Test description",
					Height:                         10.0,
					Length:                         10.0,
					NetWeight:                      100,
					ExpirationRate:                 1,
					RecommendedFreezingTemperature: 18,
					FreezingRate:                   18,
					ProductTypeID:                  1,
					SellerID:                       1,
					Width:                          10.0,
				},
			},
		},
		{
			name: "create_fail_status_422",
			mockSetup: func(m *MockProductService) {
				m.On("Create", mock.Anything).Return(internal.Product{}, internal.ErrProductUnprocessableEntity)
			},
			requestBody:      internal.Product{},
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedResponse: ResponseCreate{},
		},
		{
			name: "create_conflito_status_409",
			mockSetup: func(m *MockProductService) {
				m.On("Create", mock.Anything).Return(internal.Product{}, internal.ErrProductConflit)
			},
			requestBody:      internal.Product{},
			expectedStatus:   http.StatusConflict,
			expectedResponse: ResponseCreate{},
		},
		{
			name: "create_conflito_status_404",
			mockSetup: func(m *MockProductService) {
				m.On("Create", mock.Anything).Return(internal.Product{}, internal.ErrProductNotFound)
			},
			requestBody:      internal.Product{},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: ResponseCreate{},
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockProductService) {
				m.On("Create", mock.Anything).Return(internal.Product{}, errors.New("internal server error"))
			},
			requestBody: internal.Product{
				ProductCode:                    "Product A",
				Description:                    "descrição",
				Height:                         10.0,
				Length:                         10.0,
				NetWeight:                      100.0,
				ExpirationRate:                 1.0,
				RecommendedFreezingTemperature: 18.0,
				FreezingRate:                   18.0,
				ProductTypeID:                  1,
				SellerID:                       1,
				Width:                          10.0,
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: ResponseCreate{},
		},
		{
			name:             "create_fail_status_404",
			mockSetup:        func(m *MockProductService) {},
			requestBody:      []byte(`{invalid json}`),
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: ResponseCreate{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			productHandler := handler.NewProductHandlerDefault(mockService)
			tt.mockSetup(mockService)

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(requestBody))
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
			assert.Equal(t, actualResponse, tt.expectedResponse)

		})
	}
}

func Test_Update(t *testing.T) {
	type ResponseCreate struct {
		Data internal.Product `json:"data"`
	}
	tests := []struct {
		name             string
		mockSetup        func(*MockProductService)
		requestBody      interface{}
		id               string
		expectedStatus   int
		expectedResponse ResponseCreate
	}{
		{
			name: "update_ok_status_200",
			mockSetup: func(m *MockProductService) {
				m.On("Update", mock.Anything).Return(internal.Product{
					ID:                             1,
					ProductCode:                    "Product A",
					Description:                    "Test description",
					Height:                         10.0,
					Length:                         10.0, // Adicionado aqui se importar
					NetWeight:                      100.0,
					ExpirationRate:                 1.0,
					RecommendedFreezingTemperature: 18.0,
					FreezingRate:                   18.0,
					ProductTypeID:                  1,
					SellerID:                       1,
					Width:                          10.0,
				}, nil)
			},
			requestBody: internal.Product{
				ID:                             1,
				ProductCode:                    "Product A",
				Description:                    "Test description",
				Height:                         10.0,
				Length:                         10.0,
				NetWeight:                      100,
				ExpirationRate:                 1,
				RecommendedFreezingTemperature: 18,
				FreezingRate:                   18,
				ProductTypeID:                  1,
				SellerID:                       1,
				Width:                          10.0,
			},
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedResponse: ResponseCreate{
				Data: internal.Product{
					ID:                             1,
					ProductCode:                    "Product A",
					Description:                    "Test description",
					Height:                         10.0,
					Length:                         10.0,
					NetWeight:                      100,
					ExpirationRate:                 1,
					RecommendedFreezingTemperature: 18,
					FreezingRate:                   18,
					ProductTypeID:                  1,
					SellerID:                       1,
					Width:                          10.0,
				},
			},
		},
		{
			name: "update_fail_status_422",
			mockSetup: func(m *MockProductService) {
				m.On("Update", mock.Anything).Return(internal.Product{ID: 1}, internal.ErrProductUnprocessableEntity)
			},
			requestBody:      internal.Product{},
			id:               "1",
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedResponse: ResponseCreate{},
		},
		{
			name: "update_non_existent_status_404",
			mockSetup: func(m *MockProductService) {
				m.On("Update", mock.Anything).Return(internal.Product{ID: 1}, internal.ErrProductNotFound)
			},
			requestBody:      internal.Product{},
			id:               "1",
			expectedStatus:   http.StatusNotFound,
			expectedResponse: ResponseCreate{},
		},
		{
			name: "update_conflito_status_409",
			mockSetup: func(m *MockProductService) {
				m.On("Update", mock.Anything).Return(internal.Product{ID: 1}, internal.ErrProductConflit)
			},
			requestBody:      internal.Product{},
			id:               "1",
			expectedStatus:   http.StatusConflict,
			expectedResponse: ResponseCreate{},
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockProductService) {
				m.On("Update", mock.Anything).Return(internal.Product{ID: 1}, errors.New("internal server error"))
			},
			requestBody: internal.Product{
				ProductCode:                    "Product A",
				Description:                    "descrição",
				Height:                         10.0,
				Length:                         10.0,
				NetWeight:                      100.0,
				ExpirationRate:                 1.0,
				RecommendedFreezingTemperature: 18.0,
				FreezingRate:                   18.0,
				ProductTypeID:                  1,
				SellerID:                       1,
				Width:                          10.0,
			},
			id:               "1",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: ResponseCreate{},
		},
		{
			name: "update_fail_status_400_invalid_id",
			mockSetup: func(m *MockProductService) {
			},
			requestBody:      internal.Product{},
			id:               "invalid",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: ResponseCreate{},
		},

		{
			name: "update_fail_status_400_invalid_json",
			mockSetup: func(m *MockProductService) {
			},
			requestBody:      []byte(`{invalid json}`),
			id:               "1",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: ResponseCreate{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			productHandler := handler.NewProductHandlerDefault(mockService)
			tt.mockSetup(mockService)

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, "/product/"+tt.id, bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			productHandler.Update(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			var actualResponse ResponseCreate

			err = json.NewDecoder(rec.Body).Decode(&actualResponse)
			require.NoError(t, err)
			assert.Equal(t, actualResponse, tt.expectedResponse)

		})
	}
}

func Test_Delete(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockProductService)
		id                 string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should delete a product",
			mockSetup: func(m *MockProductService) {
				m.On("Delete", 1).Return(nil)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   nil,
		},
		{
			name: "should delete err conflit ",
			mockSetup: func(m *MockProductService) {
				m.On("Delete", 1).Return(internal.ErrProductConflit)
			},
			id:                 "1",
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   nil,
		},
		{
			name: "should return not found error",
			mockSetup: func(m *MockProductService) {
				m.On("Delete", 1).Return(internal.ErrProductNotFound)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *resterr.NewNotFoundError("product not found"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockProductService) {
				m.On("Delete", 1).Return(errors.New("internal server error"))
			},
			id:                 "1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name:               "should return bad request error",
			mockSetup:          func(m *MockProductService) {},
			id:                 "invalid_id",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			productHandler := handler.NewProductHandlerDefault(mockService)
			tt.mockSetup(mockService)

			req, err := http.NewRequest(http.MethodDelete, "/products/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			productHandler.Delete(rec, req)

			assert.Equal(t, tt.expectedStatusCode, rec.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case resterr.RestErr:
					var actualResponse resterr.RestErr
					err = json.NewDecoder(rec.Body).Decode(&actualResponse)
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

func Test_ReportRecords_All(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockProductService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "ReportRecords_All_status_200",
			mockSetup: func(p *MockProductService) {
				mockProduct := []internal.ProductRecordsJSONCount{
					{ProductID: 1, Description: "Product 1", RecordsCount: 1},
					{ProductID: 2, Description: "Product 2", RecordsCount: 1},
					{ProductID: 3, Description: "Product 3", RecordsCount: 1},
				}
				p.On("GetAllRecord").Return(mockProduct, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []internal.ProductRecordsJSONCount{
					{ProductID: 1, Description: "Product 1", RecordsCount: 1},
					{ProductID: 2, Description: "Product 2", RecordsCount: 1},
					{ProductID: 3, Description: "Product 3", RecordsCount: 1},
				},
			},
		},
		{
			name: "ReportRecords_All_status_400",
			mockSetup: func(p *MockProductService) {
				p.On("GetAllRecord").Return([]internal.ProductRecordsJSONCount{}, errors.New("bad request"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			productHandler := handler.NewProductHandlerDefault(mockService)
			tt.mockSetup(mockService)

			req := httptest.NewRequest(http.MethodGet, "/productRecords", nil)
			rec := httptest.NewRecorder()

			productHandler.ReportRecords(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != nil {
				switch response := tt.expectedBody.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data []internal.ProductRecordsJSONCount `json:"data"`
					}
					err := json.NewDecoder(rec.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)

				case resterr.RestErr:
					var actualResponse resterr.RestErr
					err := json.NewDecoder(rec.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response, actualResponse)
				default:
					t.Fatalf("Tipo de resposta inesperado: %T", response)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func Test_ReportRecords_ByID(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockProductService)
		id             string
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "ReportRecords_By_ID_status_200",
			mockSetup: func(p *MockProductService) {
				mockProduct := internal.ProductRecordsJSONCount{
					ProductID: 1, Description: "Product 1", RecordsCount: 1,
				}
				p.On("GetByIDRecord", 1).Return(mockProduct, nil)
			},
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": internal.ProductRecordsJSONCount{
					ProductID: 1, Description: "Product 1", RecordsCount: 1,
				},
			},
		},
		{
			name: "ReportRecords_By_ID_status_400",
			mockSetup: func(p *MockProductService) {
			},
			id:             "invalid_id",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   *resterr.NewBadRequestError("error parsing id with invalid syntax"),
		},
		{
			name: "ReportRecords_By_ID_status_404",
			mockSetup: func(p *MockProductService) {
				p.On("GetByIDRecord", 1).Return(internal.ProductRecordsJSONCount{}, internal.ErrProductNotFound)
			},
			id:             "1",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockProductService)
			productHandler := handler.NewProductHandlerDefault(mockService)
			tt.mockSetup(mockService)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprint("/productRecords?id=", tt.id), nil)
			rec := httptest.NewRecorder()

			productHandler.ReportRecords(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != nil {
				switch response := tt.expectedBody.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data internal.ProductRecordsJSONCount `json:"data"`
					}
					err := json.NewDecoder(rec.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					// Verifica se os dados retornados são os mesmos
					assert.Equal(t, response["data"], actualResponse.Data)

				case resterr.RestErr:
					var actualResponse resterr.RestErr
					err := json.NewDecoder(rec.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					// Verifica se o erro retornado é o mesmo
					assert.Equal(t, response, actualResponse)
				default:
					t.Fatalf("Tipo de resposta inesperado: %T", response)
				}
			}

			// Verifica se as expectativas do mock foram atendidas
			mockService.AssertExpectations(t)
		})
	}
}
