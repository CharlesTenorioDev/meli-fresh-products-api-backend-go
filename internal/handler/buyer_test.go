package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	DbPath = "../../db/buyer.json"
	Api    = "/api/v1/buyers"
)

type BuyerServiceMock struct {
	mock.Mock
}

func (bm *BuyerServiceMock) GetAll() map[int]internal.Buyer {
	args := bm.Called()
	return args.Get(0).(map[int]internal.Buyer)
}

func (bm *BuyerServiceMock) FindByID(id int) (internal.Buyer, error) {
	args := bm.Called(id)
	return args.Get(0).(internal.Buyer), args.Error(1)
}

func (bm *BuyerServiceMock) Save(newBuyer *internal.Buyer) (err error) {
	args := bm.Called(newBuyer)
	newBuyer.ID = 1
	return args.Error(0)
}

func (bm *BuyerServiceMock) Update(id int, buyerPatch internal.BuyerPatch) (err error) {
	args := bm.Called(id, buyerPatch)
	return args.Error(0)
}

func (bm *BuyerServiceMock) Delete(id int) (err error) {
	args := bm.Called(id)
	return args.Error(0)
}

func (bm *BuyerServiceMock) ReportPurchaseOrders() ([]internal.PurchaseOrdersByBuyer, error) {
	args := bm.Called()
	return args.Get(0).([]internal.PurchaseOrdersByBuyer), args.Error(1)
}

func (bm *BuyerServiceMock) ReportPurchaseOrdersByID(id int) ([]internal.PurchaseOrdersByBuyer, error) {
	args := bm.Called(id)
	return args.Get(0).([]internal.PurchaseOrdersByBuyer), args.Error(1)
}

type BuyerSuiteTest struct {
	buyersFromFile map[int]internal.Buyer
	hd             *handler.BuyerHandlerDefault
	suite.Suite
}

func (suite *BuyerSuiteTest) SetupTest() {
	suite.buyersFromFile = make(map[int]internal.Buyer)
	file, err := os.Open(DbPath)
	if err != nil {
		log.Fatal("Failed to open db file: ", err)
	}
	var buyers []internal.Buyer
	json.NewDecoder(file).Decode(&buyers) // Assume JSON is always valid

	for i, b := range buyers {
		suite.buyersFromFile[i] = b
	}
	rp, _ := repository.NewBuyerMap(DbPath)
	sv := service.NewBuyerService(rp)
	suite.hd = handler.NewBuyerHandlerDefault(sv)
}

func (suite *BuyerSuiteTest) TestGetAllBuyers() {
	r := httptest.NewRequest(http.MethodGet, Api, nil)
	w := httptest.NewRecorder()
	suite.hd.GetAll(w, r)
	assert.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)

	var buyers struct {
		Data map[string]internal.Buyer `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&buyers)
	require.NoError(suite.T(), err)

	// Validate if all buyers are correct
	for i := 0; i < len(suite.buyersFromFile); i++ {
		require.Equal(suite.T(), suite.buyersFromFile[i], buyers.Data[strconv.Itoa(i)])
	}
}

func (suite *BuyerSuiteTest) TestGetBuyersById() {
	suite.Run("get several ids", func() {
		for i := range suite.buyersFromFile {
			r := httptest.NewRequest(http.MethodGet, Api+"/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(i))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()
			suite.hd.GetByID(w, r)
			var buyers struct {
				Data internal.Buyer `json:"data"`
			}
			err := json.NewDecoder(w.Body).Decode(&buyers)
			require.NoError(suite.T(), err)
			require.Equal(suite.T(), suite.buyersFromFile[i], buyers.Data)
		}
	})
}

func (suite *BuyerSuiteTest) TestCreateBuyer() {
	// Insert the object
	bc := internal.Buyer{
		FirstName:    "Fabio",
		LastName:     "Nacarelli",
		CardNumberID: "80028922",
	}
	b, _ := json.Marshal(bc)
	r := httptest.NewRequest(http.MethodPost, Api, bytes.NewReader(b))
	w := httptest.NewRecorder()
	suite.hd.Create(w, r)
	assert.Equal(suite.T(), http.StatusCreated, w.Result().StatusCode)

	// Test if it was actually inserted
	r = httptest.NewRequest(http.MethodGet, Api+"/5", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "5")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	suite.hd.GetByID(w, r)

	var buyerCreated struct {
		Data internal.Buyer `json:"data"`
	}
	json.NewDecoder(w.Body).Decode(&buyerCreated)
	require.Equal(suite.T(), 5, buyerCreated.Data.ID)
	require.Equal(suite.T(), bc.FirstName, buyerCreated.Data.FirstName)
	require.Equal(suite.T(), bc.LastName, buyerCreated.Data.LastName)
	require.Equal(suite.T(), bc.CardNumberID, buyerCreated.Data.CardNumberID)
}

func (suite *BuyerSuiteTest) TestCreateBuyerWithMissingParameters() {
	bc := internal.Buyer{}
	b, _ := json.Marshal(bc)
	r := httptest.NewRequest(http.MethodPost, Api, bytes.NewReader(b))
	w := httptest.NewRecorder()
	suite.hd.Create(w, r)
	assert.Equal(suite.T(), http.StatusUnprocessableEntity, w.Result().StatusCode)
}

func (suite *BuyerSuiteTest) TestGetIdThatDoesntExist() {
	r := httptest.NewRequest(http.MethodGet, Api+"/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "200")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	suite.hd.GetByID(w, r)
	var buyers struct {
		Data internal.Buyer `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&buyers)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)
}

func (suite *BuyerSuiteTest) TestPatchBuyer() {
	fname := "Doe"
	lname := "John"
	cardNumberId := "404"
	suite.Run("apply changes", func() {
		bp := internal.BuyerPatch{
			FirstName:    &fname,
			LastName:     &lname,
			CardNumberID: &cardNumberId,
		}
		b, _ := json.Marshal(bp)
		r := httptest.NewRequest(http.MethodPatch, Api+"/{id}", bytes.NewReader(b))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "0")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		suite.hd.Update(w, r)
		require.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)

		var updatedBuyer struct {
			Data internal.Buyer `json:"data"`
		}
		json.NewDecoder(w.Body).Decode(&updatedBuyer)
		require.Equal(suite.T(), fname, updatedBuyer.Data.FirstName)
		require.Equal(suite.T(), lname, updatedBuyer.Data.LastName)
		require.Equal(suite.T(), cardNumberId, updatedBuyer.Data.CardNumberID)
	})

	// Check if changes were applied
	suite.Run("check if changes were applied", func() {
		r := httptest.NewRequest(http.MethodGet, Api+"/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "0")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		suite.hd.GetByID(w, r)
		var buyers struct {
			Data internal.Buyer `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&buyers)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)
		require.Equal(suite.T(), fname, buyers.Data.FirstName)
		require.Equal(suite.T(), lname, buyers.Data.LastName)
		require.Equal(suite.T(), cardNumberId, buyers.Data.CardNumberID)
	})
}

func (suite *BuyerSuiteTest) TestDeleteBuyer() {
	// Test delete
	r := httptest.NewRequest(http.MethodDelete, Api+"/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "0")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	suite.hd.Delete(w, r)
	require.Equal(suite.T(), http.StatusNoContent, w.Result().StatusCode)

	// Test if deleted
	r = httptest.NewRequest(http.MethodGet, Api+"/{id}", nil)
	rctx = chi.NewRouteContext()
	rctx.URLParams.Add("id", "0")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	suite.hd.GetByID(w, r)
	require.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)
}

func TestBuyerSuiteTest(t *testing.T) {
	suite.Run(t, new(BuyerSuiteTest))
}

// Unit tests
type TestCasesUnit struct {
	name               string
	mockService        func(*BuyerServiceMock)
	id                 string
	body               string
	expectedBody       string
	expectedStatusCode int
	expectedResponse   interface{}
	expectedMockCalls  int
}

func TestHandler_BuyerCreateUnitTest(t *testing.T) {
	testCases := []*TestCasesUnit{
		{
			name: "status code 201 (success) - Create a new buyer",
			body: `{
				"card_number_id": "123456789",
				"first_name": "Paloma",
				"last_name": "Souza"
			}`,
			expectedBody: `{"data":{"id":1,"card_number_id":"123456789","first_name":"Paloma","last_name":"Souza"}}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Save", mock.Anything).Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   internal.Buyer{},
			expectedMockCalls:  1,
		},

		{
			name: "status code 422 (fail) - Failed to create a new buyer with invalid input",
			body: `{
				"card_number_id": "",
				"first_name": "",
				"last_name": ""
			}`,
			expectedBody: `{"message":"couldn't parse buyer","error":"unprocessable_entity","code":422,"causes":null}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Save", mock.Anything).Return(errors.New("couldn't parse buyer"))
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   *resterr.NewUnprocessableEntityError("couldn't parse buyer"),
			expectedMockCalls:  1,
		},

		{
			name: "status code 409 (fail) - Failed to create a new buyer with given card number already registered",
			body: `{
				"card_number_id": "123456789",
				"first_name": "Pah",
				"last_name": "Gabi"
			}`,
			expectedBody: `{"message":"buyer with given card number already registered","error":"conflict","code":409,"causes":null}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Save", mock.Anything).Return(service.ErrCardNumberAlreadyInUse)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   *resterr.NewConflictError("buyer with given card number already registered"),
			expectedMockCalls:  1,
		},
		{
			name: "status code 409 (fail) - Failed to create a new buyer who is already registered",
			body: `{
				"card_number_id": "123456789",
				"first_name": "Paloma",
				"last_name": "Souza"
			}`,
			expectedBody: `{"message":"buyer already exists","error":"conflict","code":409,"causes":null}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Save", mock.Anything).Return(service.ErrBuyerAlreadyExists)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   *resterr.NewConflictError("buyer already exists"),
			expectedMockCalls:  1,
		},
		{
			name: "status code 400 (fail) - Attempt to create a new buyer with invalid input",
			body: `{
				"card_number_id": "8784738",
				"first_name": "Pah",
				"last_name": 1,
			}`,
			expectedBody: `{"message":"Invalid data","error":"bad_request","code":400,"causes":null}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Save", mock.Anything).Return(resterr.NewBadRequestError("invalid input data"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   *resterr.NewBadRequestError("invalid input data"),
			expectedMockCalls:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := new(BuyerServiceMock)
			hd := handler.NewBuyerHandlerDefault(sv)
			tc.mockService(sv)

			//request
			r := httptest.NewRequest(http.MethodPost, Api, strings.NewReader(tc.body))
			r.Header.Set("Content-Type", "application/json")
			//response
			w := httptest.NewRecorder()

			//When
			hd.Create(w, r)

			//Then
			require.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
			sv.AssertNumberOfCalls(t, "Save", tc.expectedMockCalls)
		})
	}

}

func TestHandler_BuyerReadUnitTest(t *testing.T) {
	testCases := []*TestCasesUnit{
		{

			name: "status code 200 (success) - Successfully read all buyer",

			expectedBody: `{"data":{"1":{"id":1,"card_number_id":"97312830","first_name":"Paloma","last_name":"Souza"},"2":{"id":2,"card_number_id":"493779","first_name":"Pah","last_name":"Gabi"}}}`,

			mockService: func(bm *BuyerServiceMock) {
				bm.On("GetAll").Return(map[int]internal.Buyer{
					1: {ID: 1, CardNumberID: "97312830", FirstName: "Paloma", LastName: "Souza"},
					2: {ID: 2, CardNumberID: "493779", FirstName: "Pah", LastName: "Gabi"},
				}, nil)
			},

			expectedStatusCode: http.StatusOK,
			expectedMockCalls:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := new(BuyerServiceMock)
			hd := handler.NewBuyerHandlerDefault(sv)
			tc.mockService(sv)

			//request
			r := httptest.NewRequest(http.MethodGet, Api, nil)
			r.Header.Set("Content-Type", "application/json")
			//response
			w := httptest.NewRecorder()

			//When
			hd.GetAll(w, r)

			//Then
			require.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
			sv.AssertNumberOfCalls(t, "GetAll", tc.expectedMockCalls)
		})
	}
}

func TestHandler_BuyerReadByIdUnitTest(t *testing.T) {
	testCases := []*TestCasesUnit{
		{

			name: "status code 404 (fail) - Attempt to get a buyer by a non existent id",
			id:   "50",

			expectedBody: `{"message":"buyer not found","error":"not_found","code":404,"causes":null}`,

			mockService: func(bm *BuyerServiceMock) {
				bm.On("FindByID", 50).Return(internal.Buyer{}, service.ErrBuyerNotFound)
			},

			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   resterr.NewNotFoundError("buyer not found"),
			expectedMockCalls:  1,
		},
		{

			name: "status code 200 (success) - Successefully get a buyer by id",
			id:   "1",

			expectedBody: `{"data":{"id":1,"card_number_id":"97312830","first_name":"Paloma","last_name":"Souza"}}`,

			mockService: func(bm *BuyerServiceMock) {
				bm.On("FindByID", 1).Return(internal.Buyer{ID: 1, CardNumberID: "97312830", FirstName: "Paloma", LastName: "Souza"}, nil)
			},

			expectedStatusCode: http.StatusOK,
			expectedMockCalls:  1,
		},
		{

			name: "status code 400 (fail) - Attempt to get a buyer with a invalid Id",
			id:   "@",

			expectedBody: `{"message":"failed to parse id","error":"bad_request","code":400,"causes":null}`,

			mockService: func(bm *BuyerServiceMock) {},

			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   *resterr.NewBadRequestError("invalid id"),
			expectedMockCalls:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := new(BuyerServiceMock)
			hd := handler.NewBuyerHandlerDefault(sv)
			tc.mockService(sv)

			//request
			r := httptest.NewRequest(http.MethodGet, Api+"/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			r.Header.Set("Content-Type", "application/json")
			//response
			w := httptest.NewRecorder()

			//When
			hd.GetByID(w, r)

			//Then
			require.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
			sv.AssertNumberOfCalls(t, "FindByID", tc.expectedMockCalls)
		})
	}
}

func TestHandler_BuyerUpdateUnitTest(t *testing.T) {
	testCases := []*TestCasesUnit{
		{
			name: "status code 200 (success) - Sucessfully update a buyer",
			id:   "0",
			body: `{
				"card_number_id": "5555555",
				"first_name": "John",
				"last_name": "Doe"
			}`,
			expectedBody: `{"data":{"card_number_id":"5555555","first_name":"John","last_name":"Doe"}}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Update", 0, mock.Anything).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedMockCalls:  1,
		},
		{
			name: "status code 404 (fail) - Attempt to update a non existent buyer",
			id:   "900",
			body: `{
				"card_number_id": "7777777",
				"first_name": "Jack",
				"last_name": "Sparrow"
			}`,
			expectedBody: `{"message":"buyer not found","error":"not_found","code":404,"causes":null}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Update", 900, mock.Anything).Return(service.ErrBuyerNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   resterr.NewNotFoundError("buyer not found"),
			expectedMockCalls:  1,
		},
		{
			name: "status code 400 (fail) - Attempt to update a buyer with invalid id",
			id:   "@",
			body: `{
				"card_number_id": "B500",
				"first_name": 123,
				"last_name": 9
			}`,
			expectedBody: `{"message":"failed to parse id","error":"bad_request","code":400,"causes":null}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Update", 0, mock.Anything).Return(resterr.NewBadRequestError("invalid id"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   *resterr.NewBadRequestError("invalid id"),
			expectedMockCalls:  0,
		},
		{
			name: "status code 400 (fail) - Attempt to update a buyer with invalid input",
			id:   "0",
			body: `{
				"card_number_id": "B500",
				"first_name": 123,
				"last_name": 9
			}`,
			expectedBody: `{"message":"failed to parse body","error":"bad_request","code":400,"causes":null}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Update", 0, mock.Anything).Return(resterr.NewBadRequestError("invalid input data"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   *resterr.NewBadRequestError("invalid input data"),
			expectedMockCalls:  0,
		},
		{
			name: "status code 409 (fail) - Failed to update a buyer with given card number already registered",
			id:   "0",
			body: `{
				"card_number_id": "123456789",
				"first_name": "John",
				"last_name": "Walker"
			}`,
			expectedBody: `{"message":"buyer with given card number already registered","error":"conflict","code":409,"causes":null}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Update", 0, mock.Anything).Return(service.ErrCardNumberAlreadyInUse)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   *resterr.NewConflictError("buyer with given card number already registered"),
			expectedMockCalls:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := new(BuyerServiceMock)
			hd := handler.NewBuyerHandlerDefault(sv)
			tc.mockService(sv)

			//request
			r := httptest.NewRequest(http.MethodPatch, Api+"/{id}", strings.NewReader(tc.body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			r.Header.Set("Content-Type", "application/json")
			//response
			w := httptest.NewRecorder()

			//When
			hd.Update(w, r)

			//Then
			require.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
			sv.AssertNumberOfCalls(t, "Update", tc.expectedMockCalls)
		})
	}

}

func TestHandler_BuyerDeleteUnitTest(t *testing.T) {
	testCases := []*TestCasesUnit{

		{
			name: "status code 404 (fail) - Attempt to delete a non existent buyer",
			id:   "670",

			expectedBody: `{"message":"buyer not found","error":"not_found","code":404,"causes":null}`,
			mockService: func(bm *BuyerServiceMock) {
				bm.On("Delete", 670).Return(service.ErrBuyerNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   resterr.NewNotFoundError("buyer not found"),
			expectedMockCalls:  1,
		},
		{
			name: "status code 204 (success) - Successfully to delete a buyer",
			id:   "3",

			mockService: func(bm *BuyerServiceMock) {
				bm.On("Delete", 3).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
			expectedMockCalls:  1,
		},
		{
			name: "status code 400 (fail) - Attempt to delete a buyer with a invalid Id",
			id:   "@",

			expectedBody: `{"message":"failed to parse id","error":"bad_request","code":400,"causes":null}`,

			mockService: func(bm *BuyerServiceMock) {},

			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   *resterr.NewBadRequestError("invalid id"),
			expectedMockCalls:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := new(BuyerServiceMock)
			hd := handler.NewBuyerHandlerDefault(sv)
			tc.mockService(sv)

			//request
			r := httptest.NewRequest(http.MethodDelete, Api+"/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			r.Header.Set("Content-Type", "application/json")
			//response
			w := httptest.NewRecorder()

			//When
			hd.Delete(w, r)

			//Then
			require.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
			sv.AssertNumberOfCalls(t, "Delete", tc.expectedMockCalls)
		})
	}
}

func TestHandler_BuyerReportPurchaseOrdersUnitTest(t *testing.T) {
	testCases := []*TestCasesUnit{

		{
			name: "status code 200 (success) - Successfully get purchase orders for all buyers",
			mockService: func(bm *BuyerServiceMock) {
				bm.On("ReportPurchaseOrders").Return([]internal.PurchaseOrdersByBuyer{
					{BuyerID: 1, CardNumberID: "123", FirstName: "Paloma", LastName: "Souza", PurchaseOrdersCount: 2},
					{BuyerID: 2, CardNumberID: "234", FirstName: "Pah", LastName: "Gabi", PurchaseOrdersCount: 1},
				}, nil)
			},

			expectedBody: `{"data":[{"id":1,"card_number_id":"123","first_name":"Paloma","last_name":"Souza","purchase_orders_count":2},{"id":2,"card_number_id":"234","first_name":"Pah","last_name":"Gabi","purchase_orders_count":1}]}`,

			expectedStatusCode: http.StatusOK,
		},
		{
			name: "status code 200 (success) - Successfully get purchase orders for a specific buyer by Id",
			id:   "2",
			mockService: func(bm *BuyerServiceMock) {
				bm.On("ReportPurchaseOrdersByID", 2).Return([]internal.PurchaseOrdersByBuyer{
					{BuyerID: 2, CardNumberID: "234", FirstName: "Pah", LastName: "Gabi", PurchaseOrdersCount: 1},
				}, nil)
			},

			expectedBody: `{"data":[{"id":2,"card_number_id":"234","first_name":"Pah","last_name":"Gabi","purchase_orders_count":1}]}`,

			expectedStatusCode: http.StatusOK,
		},
		{
			name: "status code 400 (fail) - Attempt to get purchase orders with a invalid id parameter",
			id:   "@",

			mockService: func(bm *BuyerServiceMock) {},

			expectedBody: `{"message":"failed to parse id","error":"bad_request","code":400,"causes":null}`,

			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "status code 404 (fail) - Attempt to get purchase orders for a non existent buyer by Id",
			id:   "550",

			mockService: func(bm *BuyerServiceMock) {
				bm.On("ReportPurchaseOrdersByID", 550).Return([]internal.PurchaseOrdersByBuyer{}, service.ErrBuyerNotFound)
			},

			expectedBody: `{"message":"buyer not found","error":"not_found","code":404,"causes":null}`,

			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "status code 500 (fail) - Internal server error while getting purchase orders for all buyers",

			mockService: func(bm *BuyerServiceMock) {
				bm.On("ReportPurchaseOrders").Return([]internal.PurchaseOrdersByBuyer{}, errors.New("Internal server error"))
			},

			expectedBody: `{"message":"Internal server error","error":"internal_server_error","code":500,"causes":null}`,

			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "status code 500 (fail) - Internal server error while getting purchase orders for a specific buyer by Id",
			id:   "1",

			mockService: func(bm *BuyerServiceMock) {
				bm.On("ReportPurchaseOrdersByID", 1).Return([]internal.PurchaseOrdersByBuyer{}, errors.New("Internal server error"))
			},

			expectedBody: `{"message":"Internal server error","error":"internal_server_error","code":500,"causes":null}`,

			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sv := new(BuyerServiceMock)
			hd := handler.NewBuyerHandlerDefault(sv)
			tc.mockService(sv)

			//request
			url := Api + "/report/purchase_orders"
			if tc.id != "" {
				url = Api + "/report/purchase_orders?id=" + tc.id
			}
			r := httptest.NewRequest(http.MethodGet, url, nil)
			r.Header.Set("Content-Type", "application/json")
			//response
			w := httptest.NewRecorder()

			//When
			hd.ReportPurchaseOrders(w, r)

			//Then
			require.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())

		})
	}
}
