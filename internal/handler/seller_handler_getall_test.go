package handler_test

import (
	"encoding/json"
	"errors"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/dto"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
			expectedResponse: dto.DataResponse{
				Data: []interface{}{
					map[string]interface{}{
						"id":           1,
						"cid":          123,
						"company_name": "Test Seller 1",
						"address":      "Rua 1",
						"telephone":    "1234567890",
					},
					map[string]interface{}{
						"id":           2,
						"cid":          456,
						"company_name": "Test Seller 2",
						"address":      "Rua 2",
						"telephone":    "9876543210",
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
				case dto.DataResponse:
					var actualResponse dto.DataResponse
					err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response, actualResponse)

				case rest_err.RestErr:
					var actualResponse rest_err.RestErr
					err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
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
