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

const endpoint = "/api/v1/carries"

type MockCarriesService struct {
	mock.Mock
}

func NewMockCarriesService() *MockCarriesService {
	return &MockCarriesService{}
}

func (m *MockCarriesService) FindAll() ([]internal.Carries, error) {
	args := m.Called()
	return args.Get(0).([]internal.Carries), args.Error(1)
}

func (m *MockCarriesService) Create(carry internal.Carries) (int64, error) {
	args := m.Called(carry)
	return args.Get(0).(int64), args.Error(1)
}

type UnitTestCases struct {
	name               string
	mockService        func(*MockCarriesService)
	bodyRequest        string
	expectedBody       string
	expectedStatusCode int
	expectedResponse   any
	expectedMockCalls  int
}

func TestHandler_CarriesGetAll(t *testing.T) {

	testCases := []*UnitTestCases{
		{
			name: "status code 200 (success) - Successfully read all carries",

			expectedBody: `{
				"data": [
					{
					"id": 1,
					"cid": "123",
					"company_name": "Company 1",
					"address": "Address 1",
					"phone_number": "1234567890",
					"locality_id": 1
					},
					{
					"id": 2,
					"cid": "456",
					"company_name": "Company 2",
					"address": "Address 2",
					"phone_number": "9876543210",
					"locality_id": 2
					}
				]
				}`,

			mockService: func(m *MockCarriesService) {
				m.On("FindAll").Return([]internal.Carries{
					{
						ID:          1,
						Cid:         "123",
						CompanyName: "Company 1",
						Address:     "Address 1",
						PhoneNumber: "1234567890",
						LocalityID:  1,
					},
					{
						ID:          2,
						Cid:         "456",
						CompanyName: "Company 2",
						Address:     "Address 2",
						PhoneNumber: "9876543210",
						LocalityID:  2,
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMockCalls:  1,
		},

		{
			name: "status code 500 (fail) - Failed to get all carries",

			expectedBody: `{"message":"failed to fetch carries","error":"internal_server_error","code":500,"causes":null}`,

			mockService: func(sv *MockCarriesService) {
				sv.On(("FindAll")).Return([]internal.Carries{}, errors.New("failed to fetch carries"))
			},

			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   resterr.NewInternalServerError("failed to fetch carries"),
			expectedMockCalls:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := NewMockCarriesService()
			hd := handler.NewCarriesHandlerDefault(sv)
			tc.mockService(sv)

			//request
			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			req.Header.Set("Content-Type", "application/json")
			//response
			res := httptest.NewRecorder()

			//when
			hd.GetAll(res, req)
			//then
			expected := NormalizeJSON(tc.expectedBody)
			actual := NormalizeJSON(string(res.Body.String()))
			assert.Equal(t, expected, actual)
			sv.AssertExpectations(t)
		})
	}

}

func TestHandler_CarriesCreate(t *testing.T) {

	testCases := []*UnitTestCases{

		{
			name: "status code 201 (success) - Successfully created a new carry",

			bodyRequest: `{
			"cid": "123",
			"company_name": "Company 1",
			"address": "Address 1",
			"phone_number": "1234567890",
			"locality_id": 1
		}`,

			expectedBody: `{"data":{"id":1}}`,
			mockService: func(m *MockCarriesService) {
				m.On("Create", mock.Anything).Return(int64(1), nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: internal.Carries{
				ID:          1,
				Cid:         "123",
				CompanyName: "Company 1",
				Address:     "Address 1",
				PhoneNumber: "1234567890",
				LocalityID:  1,
			},
			expectedMockCalls: 1,
		},

		{
			name: "status code 400 (fail) - Attempt to create a new carry with invalid inputs",

			bodyRequest: `{
			"cid": 4567234,
			"company_name": "Company 1",
			"address": 123,
			"phone_number": "3489891435",
			"locality_id": 1
			}`,

			expectedBody: `{"message":"failed to parse body","error":"bad_request","code":400,"causes":null}`,
			mockService: func(m *MockCarriesService) {
				m.On("Create", mock.Anything).Return(nil, errors.New("failed to parse body"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   resterr.NewBadRequestError("failed to parse body"),
			expectedMockCalls:  0,
		},
		{
			name: "status code 422 (fail) - Attempt to create a new carry with fields are missing",
			bodyRequest: `{
			"cid": "",
			"company_name": "Company 1",
			"address": "Address 123",
			"phone_number": "",
			"locality_id": 1
			}`,

			expectedBody: `{"message":"missing fields","error":"unprocessable_entity","code":422,"causes":null}`,
			mockService: func(sv *MockCarriesService) {
				sv.On("Create", mock.Anything).Return(nil, errors.New("missing fields"))
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   resterr.NewUnprocessableEntityError("missing fields"),
			expectedMockCalls:  0,
		},
		{
			name: "status code 409 (fail) - Attempt to create a new carry that already exists",
			bodyRequest: `{
			"cid": "123",
			"company_name": "Company 1234",
			"address": "Address 76543",
			"phone_number": "5668432098",
			"locality_id": 2
			}`,

			expectedBody: `{"message":"carries already exists","error":"conflict","code":409,"causes":null}`,
			mockService: func(sv *MockCarriesService) {
				sv.On("Create", mock.Anything).Return(int64(1), errors.New("carries already exists"))
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   resterr.NewConflictError("carries already exists"),
			expectedMockCalls:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := NewMockCarriesService()
			hd := handler.NewCarriesHandlerDefault(sv)
			tc.mockService(sv)

			//request
			req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(tc.bodyRequest))
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
