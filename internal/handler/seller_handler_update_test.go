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

func TestSellerDefault_Update(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSellerService)
		id                 string
		requestBody        interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should update a seller",
			mockSetup: func(m *MockSellerService) {
				mockSeller := internal.Seller{
					ID:          1,
					CID:         123,
					CompanyName: "Test Seller",
					Address:     "Rua 1",
					Telephone:   "1234567890",
				}
				m.On("FindByID", 1).Return(mockSeller, nil)
				m.On("Update", mock.Anything).Return(nil)
			},
			id: "1",
			requestBody: dto.SellersUpdateRequest{
				CompanyName: stringPtr("Updated Seller Name"),
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: dto.DataResponse{
				Data: dto.SellersGetResponse{
					Id:          1,
					Cid:         123,
					CompanyName: "Updated Seller Name",
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
			requestBody:        nil,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *rest_err.NewNotFoundError("seller not found"),
		},
		{
			name: "should return unprocessable entity error",
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
			id: "1",
			requestBody: dto.SellersUpdateRequest{
				CompanyName: stringPtr(""), // Nome da empresa inválido
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   *rest_err.NewUnprocessableEntityError("the field CompanyName is required"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockSellerService) {
				mockSeller := internal.Seller{
					ID:          1,
					CID:         123,
					CompanyName: "Test Seller",
					Address:     "Rua 1",
					Telephone:   "1234567890",
				}
				m.On("FindByID", 1).Return(mockSeller, nil)
				m.On("Update", mock.Anything).Return(errors.New("internal server error"))
			},
			id: "1",
			requestBody: dto.SellersUpdateRequest{
				CompanyName: stringPtr("Updated Seller Name"),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name:               "should return bad request error",
			mockSetup:          func(m *MockSellerService) {},
			id:                 "invalid_id",
			requestBody:        nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSellerService)
			sellerHandler := handler.NewSellerDefault(mockService)
			tt.mockSetup(mockService)

			var requestBody []byte
			if tt.requestBody != nil {
				var err error
				requestBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			req, err := http.NewRequest(http.MethodPatch, "/sellers/"+tt.id, bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("content-type", "application/json")

			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			hd := sellerHandler.Update()
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
					assert.Equal(t, response.Data, actualResponse.Data)
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
