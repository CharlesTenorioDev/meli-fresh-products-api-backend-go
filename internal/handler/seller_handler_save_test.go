package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
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
			requestBody: dto.SellersPostRequest{
				CID:         123,
				CompanyName: "Test Seller",
				Address:     "Rua 1",
				Telephone:   "1234567890",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: dto.DataResponse{
				Data: dto.SellerPostResponse{
					Id: 1,
				},
			},
		},
		{
			name: "should return unprocessable entity error",
			mockSetup: func(m *MockSellerService) {
			},
			requestBody: dto.SellersPostRequest{
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
			requestBody: dto.SellersPostRequest{
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
			requestBody: dto.SellersPostRequest{
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
				case dto.DataResponse:
					var actualResponse dto.DataResponse
					err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response, actualResponse.Data)
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
