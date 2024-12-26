package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/dto"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockSellerService é uma estrutura mock para internal.SellerService
type MockSellerService struct {
	mock.Mock
}

// FindAll mock
func (m *MockSellerService) FindAll() ([]internal.Seller, error) {
	args := m.Called()
	return args.Get(0).([]internal.Seller), args.Error(1)
}

// FindByID mock
func (m *MockSellerService) FindByID(id int) (internal.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Seller), args.Error(1)
}

// Save mock
func (m *MockSellerService) Save(seller *internal.Seller) error {
	args := m.Called(seller)
	seller.ID = 1
	return args.Error(0)
}

// Update mock
func (m *MockSellerService) Update(id int, updatedSeller internal.SellerPatch) (internal.Seller, error) {
	args := m.Called(id, updatedSeller)
	return args.Get(1).(internal.Seller), args.Error(0)
}

// Delete mock
func (m *MockSellerService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func TestSellerDefault_GetAll(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSellerService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should return all sellers",
			mockSetup: func(m *MockSellerService) {
				mockSellers := []internal.Seller{
					{
						ID:          1,
						CID:         123,
						CompanyName: "Test Seller 1",
						Address:     "Rua 1",
						Telephone:   "1234567890",
					},
					{
						ID:          2,
						CID:         456,
						CompanyName: "Test Seller 2",
						Address:     "Rua 2",
						Telephone:   "9876543210",
					},
				}
				m.On("FindAll").Return(mockSellers, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": []dto.SellersGetJson{
					{
						Id:          1,
						Cid:         123,
						CompanyName: "Test Seller 1",
						Address:     "Rua 1",
						Telephone:   "1234567890",
					},
					{
						Id:          2,
						Cid:         456,
						CompanyName: "Test Seller 2",
						Address:     "Rua 2",
						Telephone:   "9876543210",
					},
				},
			}},
		{
			name: "should return not found error",
			mockSetup: func(m *MockSellerService) {
				m.On("FindAll").Return([]internal.Seller{}, internal.ErrSellerNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *rest_err.NewNotFoundError("sellers not found"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockSellerService) {
				m.On("FindAll").Return([]internal.Seller{}, errors.New("internal server error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil, // Não há resposta esperada específica para erro interno
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSellerService)
			sellerHandler := handler.NewSellerDefault(mockService)
			tt.mockSetup(mockService)

			req, err := http.NewRequest(http.MethodGet, "/sellers", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			hd := sellerHandler.GetAll()
			hd.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data []dto.SellersGetJson `json:"data"`
					}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)

				case rest_err.RestErr:
					var actualResponse rest_err.RestErr
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

func TestSellerDefault_GetByID(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSellerService)
		id                 string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should return a seller by ID",
			mockSetup: func(m *MockSellerService) {
				mockSeller := internal.Seller{
					ID:          1,
					CID:         123,
					CompanyName: "Test Seller",
					Address:     "Rua 1",
					Telephone:   "1234567890",
				}
				m.On("FindByID", 1).Return(mockSeller, nil)
			},
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": dto.SellersGetJson{
					Id:          1,
					Cid:         123,
					CompanyName: "Test Seller",
					Address:     "Rua 1",
					Telephone:   "1234567890",
				},
			},
		},
		{
			name: "should return not found error",
			mockSetup: func(m *MockSellerService) {
				m.On("FindByID", 1).Return(internal.Seller{}, internal.ErrSellerNotFound)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *rest_err.NewNotFoundError("seller not found"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockSellerService) {
				m.On("FindByID", 1).Return(internal.Seller{}, errors.New("internal server error"))
			},
			id:                 "1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name:               "should return bad request error",
			mockSetup:          func(m *MockSellerService) {},
			id:                 "invalid_id",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSellerService)
			sellerHandler := handler.NewSellerDefault(mockService)
			tt.mockSetup(mockService)

			req, err := http.NewRequest(http.MethodGet, "/sellers/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			hd := sellerHandler.GetByID()
			hd.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data dto.SellersGetJson `json:"data"`
					}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)
				case rest_err.RestErr:
					var actualResponse rest_err.RestErr
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

func TestSellerDefault_Save(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSellerService)
		requestBody        interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should create a new seller",
			mockSetup: func(m *MockSellerService) {
				m.On("Save", mock.Anything).Return(nil)
			},
			requestBody: dto.SellersPostJson{
				CID:         123,
				CompanyName: "Test Seller",
				Address:     "Rua 1",
				Telephone:   "1234567890",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: map[string]interface{}{
				"data": map[string]int{
					"seller_id": 1,
				},
			},
		},
		{
			name: "should return unprocessable entity error",
			mockSetup: func(m *MockSellerService) {
			},
			requestBody: dto.SellersPostJson{
				CID:         123,
				CompanyName: "", // Nome da empresa inválido
				Address:     "Rua 1",
				Telephone:   "1234567890",
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   *rest_err.NewUnprocessableEntityError("seller.CompanyName is required"),
		},
		{
			name: "should return conflict error",
			mockSetup: func(m *MockSellerService) {
				m.On("Save", mock.Anything).Return(internal.ErrSellerConflict)
			},
			requestBody: dto.SellersPostJson{
				CID:         123,
				CompanyName: "Test Seller",
				Address:     "Rua 1",
				Telephone:   "1234567890",
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   *rest_err.NewConflictError("seller already exists"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockSellerService) {
				m.On("Save", mock.Anything).Return(errors.New("internal server error"))
			},
			requestBody: dto.SellersPostJson{
				CID:         123,
				CompanyName: "Test Seller",
				Address:     "Rua 1",
				Telephone:   "1234567890",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSellerService)
			sellerHandler := handler.NewSellerDefault(mockService)
			tt.mockSetup(mockService)

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/sellers", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("content-type", "application/json")

			rr := httptest.NewRecorder()
			hd := sellerHandler.Save()

			hd.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse = struct {
						Data map[string]int `json:"data"`
					}{}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)
				case rest_err.RestErr:
					var actualResponse rest_err.RestErr
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

func TestSellerDefault_Delete(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSellerService)
		id                 string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should delete a seller",
			mockSetup: func(m *MockSellerService) {
				m.On("Delete", 1).Return(nil)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   nil,
		},
		{
			name: "should return not found error",
			mockSetup: func(m *MockSellerService) {
				m.On("Delete", 1).Return(internal.ErrSellerNotFound)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *rest_err.NewNotFoundError("seller not found"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockSellerService) {
				m.On("Delete", 1).Return(errors.New("internal server error"))
			},
			id:                 "1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name:               "should return bad request error",
			mockSetup:          func(m *MockSellerService) {},
			id:                 "invalid_id",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSellerService)
			sellerHandler := handler.NewSellerDefault(mockService)
			tt.mockSetup(mockService)

			req, err := http.NewRequest(http.MethodDelete, "/sellers/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			hd := sellerHandler.Delete()
			hd.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case rest_err.RestErr:
					var actualResponse rest_err.RestErr
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
