package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func NewPurchaseOrderMock() *PurchaseOrderServiceMock {
	return &PurchaseOrderServiceMock{}
}

type PurchaseOrderServiceMock struct {
	mock.Mock
}

func (m *PurchaseOrderServiceMock) FindByID(id int) (internal.PurchaseOrder, error) {
	args := m.Called(id)
	return args.Get(0).(internal.PurchaseOrder), args.Error(1)
}

func (m *PurchaseOrderServiceMock) Save(p *internal.PurchaseOrder) error {
	args := m.Called(p)
	return args.Error(0)
}

var (
	endpointPurchaseOrder = "/api/v1/purchase-orders"
)

func TestPurchaseOrder_Create(t *testing.T) {
	testCases := []struct {
		description       string
		method            string
		url               string
		body              string
		expectedBody      string
		expectedCode      int
		expectedHeader    http.Header
		mock              func() *PurchaseOrderServiceMock
		expectedMockCalls int
	}{
		{
			description: "case 1 - success: Create a new purchase order",
			method:      http.MethodPost,
			url:         endpointPurchaseOrder,
			body: `{
				"order_number": "123456",
				"order_date": "2021-01-01",
				"tracking_code": "123456",
				"buyer_id": 1,
				"product_record_id": 1
			}`,
			expectedBody:   `{"data":{"id":1,"order_number":"123456","order_date":"2021-01-01","tracking_code":"123456","buyer_id":1,"product_record_id":1}}`,
			expectedCode:   http.StatusCreated,
			expectedHeader: jsonHeader,
			mock: func() *PurchaseOrderServiceMock {
				mk := NewPurchaseOrderMock()
				mk.On("Save", mock.AnythingOfType("*internal.PurchaseOrder")).Run(func(args mock.Arguments) {
					w := args.Get(0).(*internal.PurchaseOrder)
					w.ID = 1
				}).Return(nil)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 2 - error: Attempt to create a purchase order without required fields",
			method:      http.MethodPost,
			url:         endpointPurchaseOrder,
			body: `{
			}`,
			expectedBody: `
				{
					"message": "purchase order inputs are missing",
					"error": "unprocessable_entity",
					"code": 422,
					"causes": [
						{
							"field": "order_number",
							"message": "order number is required"
						},
						{
							"field": "order_date",
							"message": "order date is required"
						},
						{
							"field": "tracking_code",
							"message": "tracking code is required"
						},
						{
							"field": "buyer_id",
							"message": "buyer id is required"
						},
						{
							"field": "product_record_id",
							"message": "product record id is required"
						}
					]
				}
			`,
			expectedCode:   http.StatusUnprocessableEntity,
			expectedHeader: jsonHeader,
			mock: func() *PurchaseOrderServiceMock {
				return NewPurchaseOrderMock()
			},
			expectedMockCalls: 0,
		},
		{
			description: "case 3 - error: Attempt to create a purchase order with invalid date format",
			method:      http.MethodPost,
			url:         endpointPurchaseOrder,
			body: `{
				"order_number": "123ABC",
				"order_date": "2023-10-05invalid",
				"tracking_code": "123456",
				"buyer_id": 1,
				"product_record_id": 1
			}`,
			expectedBody: `
				{
					"message": "Invalid data",
					"error": "bad_request",
					"code": 400,
					"causes": [
						{
							"field": "order_date",
							"message": "invalid date format"
						}
					]
				}
			`,
			expectedCode:   http.StatusBadRequest,
			expectedHeader: jsonHeader,
			mock: func() *PurchaseOrderServiceMock {
				return NewPurchaseOrderMock()
			},
			expectedMockCalls: 0,
		},
		{
			description: "case 4 - error: Attempt to create a purchase order with a non-existent buyer",
			method:      http.MethodPost,
			url:         endpointPurchaseOrder,
			body: `{
				"order_number": "123ABC",
				"order_date": "2023-10-05",
				"tracking_code": "123456",
				"buyer_id": 10,
				"product_record_id": 1
			}`,
			expectedBody: `
				{
					"message": "buyer not found",
					"error": "conflict",
					"code": 409,
					"causes": null
				}
			`,
			expectedCode:   http.StatusConflict,
			expectedHeader: jsonHeader,
			mock: func() *PurchaseOrderServiceMock {
				mk := NewPurchaseOrderMock()
				mk.On("Save", mock.Anything).Return(service.ErrBuyerNotFound)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 5 - error: Attempt to create a purchase order with a non-existent product record",
			method:      http.MethodPost,
			url:         endpointPurchaseOrder,
			body: `{
				"order_number": "123ABC",
				"order_date": "2023-10-05",
				"tracking_code": "123456",
				"buyer_id": 1,
				"product_record_id": 10
			}`,
			expectedBody: `
				{
					"message": "product-records not found",
					"error": "conflict",
					"code": 409,
					"causes": null
				}
			`,
			expectedCode:   http.StatusConflict,
			expectedHeader: jsonHeader,
			mock: func() *PurchaseOrderServiceMock {
				mk := NewPurchaseOrderMock()
				mk.On("Save", mock.Anything).Return(internal.ErrProductRecordsNotFound)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 6 - error: Attempt to create a new purchase order generating an unexpected error",
			method:      http.MethodPost,
			url:         endpointPurchaseOrder,
			body: `{
				"order_number": "123ABC",
				"order_date": "2023-10-05",
				"tracking_code": "123456",
				"buyer_id": 1,
				"product_record_id": 1
			}`,
			expectedBody: `
				{
					"message": "Internal Server Error",
					"error": "internal_server_error",
					"code": 500,
					"causes": null
				}
			`,
			expectedCode:   http.StatusInternalServerError,
			expectedHeader: jsonHeader,
			mock: func() *PurchaseOrderServiceMock {
				mk := NewPurchaseOrderMock()
				mk.On("Save", mock.Anything).Return(errors.New("unexpected error"))
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 7 - error: Attempt to create a new purchase order with an existing order number",
			method:      http.MethodPost,
			url:         endpointPurchaseOrder,
			body: `{
				"order_number": "123456",
				"order_date": "2023-10-05",
				"tracking_code": "123456",
				"buyer_id": 1,
				"product_record_id": 1
			}`,
			expectedBody: `
				{
					"message": "purchase order number already exists",
					"error": "conflict",
					"code": 409,
					"causes": null
				}
			`,
			expectedCode:   http.StatusConflict,
			expectedHeader: jsonHeader,
			mock: func() *PurchaseOrderServiceMock {
				mk := NewPurchaseOrderMock()
				mk.On("Save", mock.Anything).Return(internal.ErrPurchaseOrderConflict)
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 8 - error: Attempt to create a new purchase order with a negative buyer ID",
			method:      http.MethodPost,
			url:         endpointPurchaseOrder,
			body: `{
				"order_number": "123456",
				"order_date": "2023-10-05",
				"tracking_code": "123456",
				"buyer_id": -1,
				"product_record_id": 1
			}`,
			expectedBody: `
				{
					"message": "purchase order inputs are invalid",
					"error": "bad_request",
					"code": 400,
					"causes": [
						{
							"field": "buyer_id",
							"message": "buyer ID cannot be negative"
						}
					]
				}
			`,
			expectedCode:   http.StatusBadRequest,
			expectedHeader: jsonHeader,
			mock: func() *PurchaseOrderServiceMock {
				mk := NewPurchaseOrderMock()
				mk.On("Save", mock.Anything).Return(internal.DomainError{
					Message: internal.ErrPurchaseOrderBadRequest.Error(),
					Causes: []internal.Causes{
						{
							Field:   "buyer_id",
							Message: "buyer ID cannot be negative",
						},
					},
				})
				return mk
			},
			expectedMockCalls: 1,
		},
		{
			description: "case 9 - error: Attempt to create a new purchase order with invalid type of order number",
			method:      http.MethodPost,
			url:         endpointPurchaseOrder,
			body: `{
				"order_number": 123456,
				"order_date": "2023-10-05",
				"tracking_code": "123456",
				"buyer_id": 1,
				"product_record_id": 1
			}`,
			expectedBody: `
				{
					"message": "Invalid data",
					"error": "bad_request",
					"code": 400,
					"causes": null
				}
			`,
			expectedCode:   http.StatusBadRequest,
			expectedHeader: jsonHeader,
			mock: func() *PurchaseOrderServiceMock {
				return NewPurchaseOrderMock()
			},
			expectedMockCalls: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// handler and its dependencies
			sv := tc.mock()
			hd := handler.NewPurchaseOrderHandler(sv)
			hdFunc := hd.Create()

			// http request and response
			request := httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			// WHEN
			hdFunc(response, request)

			// THEN
			require.Equal(t, tc.expectedCode, response.Code)
			require.JSONEq(t, tc.expectedBody, response.Body.String())
			sv.AssertNumberOfCalls(t, "Save", tc.expectedMockCalls)
		})
	}
}
