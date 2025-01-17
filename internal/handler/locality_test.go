package handler_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	name = "txdb_localities"
)

// MockLocalityService é uma estrutura mock para internal.LocalityService
type MockLocalityService struct {
	mock.Mock
}

func (m *MockLocalityService) ReportCarries(localityId int) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockLocalityService) GetAmountOfCarriesForEveryLocality() ([]internal.CarriesCountPerLocality, error) {
	//TODO implement me
	panic("implement me")
}

// ReportSellers mock
func (m *MockLocalityService) ReportSellers() ([]internal.Locality, error) {
	args := m.Called()
	return args.Get(0).([]internal.Locality), args.Error(1)
}

// ReportSellersByID mock
func (m *MockLocalityService) ReportSellersByID(id int) ([]internal.Locality, error) {
	args := m.Called(id)
	return args.Get(0).([]internal.Locality), args.Error(1)
}

// Save mock
func (m *MockLocalityService) Save(locality *internal.Locality) error {
	args := m.Called(locality)
	return args.Error(0)
}

func (m *MockLocalityService) FindByID(id int) (locality internal.Locality, err error) {
	args := m.Called(id)
	return args.Get(0).(internal.Locality), args.Error(1)
}

type LocalityTestSuite struct {
	db *sql.DB
	hd *handler.LocalityDefault
	suite.Suite
}

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "meli_pass",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "melifresh",
	}
	txdb.Register(name, "mysql", cfg.FormatDSN())
}

func (l *LocalityTestSuite) SetupTest() {
	var err error
	l.db, err = sql.Open(name, "")
	require.NoError(l.T(), err)
	rp := repository.NewLocalityMysql(l.db)
	sv := service.NewLocalityDefault(rp)
	l.hd = handler.NewLocalityDefault(sv)
}

func (l *LocalityTestSuite) TestLocalityDefault_ReportCarries() {
	defer l.db.Close()
	l.T().Run("carries with locality id 3", func(t *testing.T) {
		expectedAmountOfCarries := 1
		r := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportCarries?id=3", nil)
		w := httptest.NewRecorder()
		l.hd.ReportCarries()(w, r)
		var res struct {
			Data struct {
				AmountOfCarries int `json:"amount_of_carries"`
			} `json:"data"`
		}
		json.NewDecoder(w.Result().Body).Decode(&res)
		require.Equal(t, expectedAmountOfCarries, res.Data.AmountOfCarries)
	})
	l.T().Run("carries without passing locality", func(t *testing.T) {
		expectedCarriesCountPerLocality := []internal.CarriesCountPerLocality{
			{
				CarriesCount: 1,
				LocalityId:   1,
				LocalityName: "New York City",
			},
			{
				CarriesCount: 1,
				LocalityId:   2,
				LocalityName: "Los Angeles",
			},
			{
				CarriesCount: 1,
				LocalityId:   3,
				LocalityName: "Chicago",
			},
			{
				CarriesCount: 1,
				LocalityId:   4,
				LocalityName: "Houston",
			},
			{
				CarriesCount: 1,
				LocalityId:   5,
				LocalityName: "Phoenix",
			},
		}
		r := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportCarries", nil)
		w := httptest.NewRecorder()
		l.hd.ReportCarries()(w, r)
		var res struct {
			Data []internal.CarriesCountPerLocality `json:"data"`
		}
		json.NewDecoder(w.Result().Body).Decode(&res)
		require.Equal(t, expectedCarriesCountPerLocality, res.Data)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
	l.T().Run("carries with a location that does not exist", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportCarries?id=100", nil)
		w := httptest.NewRecorder()
		l.hd.ReportCarries()(w, r)
		require.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func TestLocalityTestSuite(t *testing.T) {
	suite.Run(t, new(LocalityTestSuite))
}

func TestLocalityDefault_ReportSellers(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockLocalityService)
		id                 string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should return all localities with sellers count",
			mockSetup: func(m *MockLocalityService) {
				mockLocalities := []internal.Locality{
					{
						ID:           1,
						LocalityName: "Test Locality 1",
						ProvinceName: "Test Province 1",
						CountryName:  "Test Country 1",
						Sellers:      10,
					},
					{
						ID:           2,
						LocalityName: "Test Locality 2",
						ProvinceName: "Test Province 2",
						CountryName:  "Test Country 2",
						Sellers:      5,
					},
				}
				m.On("ReportSellers").Return(mockLocalities, nil)
			},
			id:                 "",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": []handler.LocalityGetJson{
					{
						ID:           1,
						LocalityName: "Test Locality 1",
						ProvinceName: "Test Province 1",
						CountryName:  "Test Country 1",
						SellersCount: 10,
					},
					{
						ID:           2,
						LocalityName: "Test Locality 2",
						ProvinceName: "Test Province 2",
						CountryName:  "Test Country 2",
						SellersCount: 5,
					},
				},
			},
		},
		{
			name: "should return a locality with sellers count by ID",
			mockSetup: func(m *MockLocalityService) {
				mockLocalities := []internal.Locality{
					{
						ID:           1,
						LocalityName: "Test Locality 1",
						ProvinceName: "Test Province 1",
						CountryName:  "Test Country 1",
						Sellers:      10,
					},
				}
				m.On("ReportSellersByID", 1).Return(mockLocalities, nil)
			},
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": []handler.LocalityGetJson{
					{
						ID:           1,
						LocalityName: "Test Locality 1",
						ProvinceName: "Test Province 1",
						CountryName:  "Test Country 1",
						SellersCount: 10,
					},
				},
			},
		},
		{
			name: "should return not found error",
			mockSetup: func(m *MockLocalityService) {
				m.On("ReportSellers").Return([]internal.Locality{}, internal.ErrLocalityNotFound)
			},
			id:                 "",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *rest_err.NewNotFoundError("locality not found"),
		},
		{
			name: "should return not found error by ID",
			mockSetup: func(m *MockLocalityService) {
				m.On("ReportSellersByID", 1).Return([]internal.Locality{}, internal.ErrLocalityNotFound)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *rest_err.NewNotFoundError("locality not found"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockLocalityService) {
				m.On("ReportSellers").Return([]internal.Locality{}, errors.New("internal server error"))
			},
			id:                 "",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name: "should return internal server error by ID",
			mockSetup: func(m *MockLocalityService) {
				m.On("ReportSellersByID", 1).Return([]internal.Locality{}, errors.New("internal server error"))
			},
			id:                 "1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name:               "should return bad request error",
			mockSetup:          func(m *MockLocalityService) {},
			id:                 "invalid_id",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   *rest_err.NewBadRequestError("id should be a number"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockLocalityService)
			localityHandler := handler.NewLocalityDefault(mockService)
			tt.mockSetup(mockService)

			req, err := http.NewRequest(http.MethodGet, "/localities/reportSellers?id="+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			hd := localityHandler.ReportSellers()
			hd.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data []handler.LocalityGetJson `json:"data"`
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

func TestLocalityDefault_Save(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockLocalityService)
		requestBody        interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should create a new locality",
			mockSetup: func(m *MockLocalityService) {
				m.On("Save", mock.Anything).Return(nil)
			},
			requestBody: handler.LocalityPostJson{
				LocalityID:   123,
				LocalityName: "Test Locality",
				ProvinceName: "Test Province",
				CountryName:  "Test Country",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": handler.LocalityPostJson{
					LocalityID:   123,
					LocalityName: "Test Locality",
					ProvinceName: "Test Province",
					CountryName:  "Test Country",
				},
			},
		},
		{
			name: "should return conflict error",
			mockSetup: func(m *MockLocalityService) {
				m.On("Save", mock.Anything).Return(internal.ErrLocalityConflict)
			},
			requestBody: handler.LocalityPostJson{
				LocalityID:   123,
				LocalityName: "Test Locality",
				ProvinceName: "Test Province",
				CountryName:  "Test Country",
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   *rest_err.NewConflictError("locality conflict"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockLocalityService) {
				m.On("Save", mock.Anything).Return(errors.New("internal server error"))
			},
			requestBody: handler.LocalityPostJson{
				LocalityID:   123,
				LocalityName: "Test Locality",
				ProvinceName: "Test Province",
				CountryName:  "Test Country",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name: "should return bad request error",
			mockSetup: func(m *MockLocalityService) {
				m.On("Save", mock.Anything).Return(internal.DomainError{
					Message: "locality validation error",
					Causes: []internal.Causes{
						{Field: "LocalityName", Message: "locality name is required"},
					},
				})
			},
			requestBody: handler.LocalityPostJson{
				LocalityID:   123,
				LocalityName: "",
				ProvinceName: "Test Province",
				CountryName:  "Test Country",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: *rest_err.NewBadRequestValidationError("locality validation error", []rest_err.Causes{
				{Field: "LocalityName", Message: "locality name is required"},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockLocalityService)
			localityHandler := handler.NewLocalityDefault(mockService)
			tt.mockSetup(mockService)

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/localities", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("content-type", "application/json")

			rr := httptest.NewRecorder()
			hd := localityHandler.Save()

			hd.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse = struct {
						Data handler.LocalityPostJson `json:"data"`
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
