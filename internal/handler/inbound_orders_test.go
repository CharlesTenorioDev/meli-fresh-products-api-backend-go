package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type InboundOrdersServiceMock struct {
	mock.Mock
}

func (inb *InboundOrdersServiceMock) Create(internal.InboundOrders) (int64, error) {
	args := inb.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (inb *InboundOrdersServiceMock) FindAll() ([]internal.InboundOrders, error) {
	args := inb.Called()
	return args.Get(0).([]internal.InboundOrders), args.Error(1)

}

type TestUnitCases struct {
	name               string
	mockService        func(*InboundOrdersServiceMock)
	body               string
	expectedBody       string
	expectedStatusCode int
	expectedResponse   interface{}
	expectedMockCalls  int
}

func TestInboundCreate(t *testing.T) {
	testCases := []*TestUnitCases{
		{
			name: "status code 201 (sucess) - Create a new Inbound Order",
			body: `{
				"order_date": "2123-01-01",
				"order_number": "ORD11111",
				"employee_id": 1,
				"product_batch_id": 1,
				"warehouse_id": 1
			}`,
			expectedBody: `{"data":{"id":1}}`,
			mockService: func(inb *InboundOrdersServiceMock) {
				inb.On("Create", mock.Anything).Return(int64(1), nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   internal.InboundOrders{},
			expectedMockCalls:  1,
		},
		{
			name: "status code 400 (fail) - Attempt to create a new Inbound Order with invalid input",
			body: `{
				"order_date": "2021-03-04",
				"order_number": "ORD2222",
				"employee_id": 1,
				"product_batch_id": "aaaaa",
				"warehouse_id": "bbbbbb"
			}`,
			expectedBody: `{"error":"invalid body format"}`,
			mockService: func(inb *InboundOrdersServiceMock) {
				inb.On("Create", mock.Anything).Return(nil, errors.New("invalid body format"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   errors.New("invalid body format"),
			expectedMockCalls:  0,
		},
		{
			name: "status code 422 (fail) - Attempt to create a new Inbound Order with fields are missing",
			body: `{
				"order_date": "2021-03-04",
				"order_number": "",
				"employee_id": 1,
				"product_batch_id": 3,
				"warehouse_id": 0
			}`,
			expectedBody: `{"error":"required fields are missing"}`,
			mockService: func(inb *InboundOrdersServiceMock) {
				inb.On("Create", mock.Anything).Return(nil, errors.New("required fields are missing"))
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   errors.New("required fields are missing"),
			expectedMockCalls:  0,
		},
		{
			name: "status code 409 (fail) - Attempt to create a new Inbound Order that already exists",
			body: `{
				"order_date": "2021-03-04",
				"order_number": "ORD11111",
				"employee_id": 1,
				"product_batch_id": 3,
				"warehouse_id": 2
			}`,
			expectedBody: `{"error":"order number already exists"}`,
			mockService: func(inb *InboundOrdersServiceMock) {
				inb.On("Create", mock.Anything).Return(int64(1), internal.ErrOrderNumberAlreadyExists)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   errors.New("order number already exists"),
			expectedMockCalls:  1,
		},
		{
			name: "status code 404 (fail) - Attempt to create a new Inbound Order with Id employee not exists",
			body: `{
				"order_date": "2021-03-04",
				"order_number": "88080",
				"employee_id": 155,
				"product_batch_id": 3,
				"warehouse_id": 2
			}`,
			expectedBody: `{"error":"employee not exists"}`,
			mockService: func(inb *InboundOrdersServiceMock) {
				inb.On("Create", mock.Anything).Return(int64(1), internal.ErrEmployeeNotFound)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   errors.New("employee not exists"),
			expectedMockCalls:  1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := new(InboundOrdersServiceMock)
			hd := handler.NewInboundOrdersHandler(sv)
			tc.mockService(sv)

			//request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/inbound-orders", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			//response
			res := httptest.NewRecorder()

			//when
			hd.Create(res, req)

			//then
			require.Equal(t, tc.expectedStatusCode, res.Code)
			assert.Equal(t, tc.expectedBody, res.Body.String())
			sv.AssertNumberOfCalls(t, "Create", tc.expectedMockCalls)
		})
	}
}

func TestInboundGetAll(t *testing.T) {
	testCases := []*TestUnitCases{

		{
			name: "status code 200 (success) - Successfully get all Inbound Orders",

			expectedBody: `{"data":[{"id":1,"order_date":"2021-03-04","order_number":"ORD11111","employee_id":1,"product_batch_id":3,"warehouse_id":2},{"id":2,"order_date":"2021-03-04","order_number":"ORD2222","employee_id":1,"product_batch_id":3,"warehouse_id":2}]}`,

			mockService: func(inb *InboundOrdersServiceMock) {
				inb.On(("FindAll")).Return([]internal.InboundOrders{
					{
						ID:             1,
						OrderDate:      "2021-03-04",
						OrderNumber:    "ORD11111",
						EmployeeID:     1,
						ProductBatchID: 3,
						WarehouseID:    2,
					},
					{
						ID:             2,
						OrderDate:      "2021-03-04",
						OrderNumber:    "ORD2222",
						EmployeeID:     1,
						ProductBatchID: 3,
						WarehouseID:    2,
					},
				}, nil)
			},

			expectedStatusCode: http.StatusOK,

			expectedMockCalls: 1,
		},

		{
			name: "status code 500 (fail) - Failed to get all Inbound Orders",

			expectedBody: `{"message":"failed to fetch inbounds orders","error":"internal_server_error","code":500,"causes":null}`,

			mockService: func(inb *InboundOrdersServiceMock) {
				inb.On(("FindAll")).Return([]internal.InboundOrders{}, errors.New("failed to fetch inbounds orders"))
			},

			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   resterr.NewInternalServerError("failed to fetch inbounds orders"),
			expectedMockCalls:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := new(InboundOrdersServiceMock)
			hd := handler.NewInboundOrdersHandler(sv)
			tc.mockService(sv)

			//request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/inbound-orders", nil)
			req.Header.Set("Content-Type", "application/json")
			//response
			res := httptest.NewRecorder()

			//when
			hd.GetAll(res, req)

			//then
			require.Equal(t, tc.expectedStatusCode, res.Code)
			assert.Equal(t, tc.expectedBody, res.Body.String())
			sv.AssertNumberOfCalls(t, "FindAll", tc.expectedMockCalls)
		})
	}
}
