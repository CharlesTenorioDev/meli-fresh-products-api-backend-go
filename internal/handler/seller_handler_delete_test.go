package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
