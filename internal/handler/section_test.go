package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSectionService struct {
	mock.Mock
}

type RequestSectionErrorJSON struct {
	SectionNumber      int     `json:"section_number"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseID        int     `json:"warehouse_id"`
	ProductTypeID      int     `json:"product_type_id"`
}

func (m *MockSectionService) FindAll() ([]internal.Section, error) {
	args := m.Called()
	return args.Get(0).([]internal.Section), args.Error(1)
}

func (m *MockSectionService) FindByID(id int) (internal.Section, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (m *MockSectionService) ReportProducts() ([]internal.ReportProduct, error) {
	args := m.Called()
	return args.Get(0).([]internal.ReportProduct), args.Error(1)
}

func (m *MockSectionService) ReportProductsByID(id int) (internal.ReportProduct, error) {
	args := m.Called(id)
	return args.Get(0).(internal.ReportProduct), args.Error(1)
}

func (m *MockSectionService) Save(section *internal.Section) error {
	if ok := section.Ok(); !ok {
		return internal.ErrSectionUnprocessableEntity
	}

	args := m.Called(section)
	section.ID = 1
	return args.Error(0)
}

func (m *MockSectionService) Update(id int, body internal.SectionPatch) (internal.Section, error) {
	args := m.Called(id, body)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (m *MockSectionService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestHandler_CreateSectionUnitTest(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSectionService)
		requestBody        interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "successfully create a new section",
			mockSetup: func(m *MockSectionService) {
				m.On("Save", mock.Anything).Return(nil)
			},
			requestBody: handler.RequestSectionJSON{
				SectionNumber:      123,
				CurrentTemperature: 22.5,
				MinimumTemperature: 15.0,
				CurrentCapacity:    50,
				MinimumCapacity:    30,
				MaximumCapacity:    100,
				WarehouseID:        1,
				ProductTypeID:      2,
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: internal.Section{
				ID:                 1,
				SectionNumber:      123,
				CurrentTemperature: 22.5,
				MinimumTemperature: 15.0,
				CurrentCapacity:    50,
				MinimumCapacity:    30,
				MaximumCapacity:    100,
				WarehouseID:        1,
				ProductTypeID:      2,
			},
		},
		{
			name: "erro campo obrigatorio",
			mockSetup: func(m *MockSectionService) {
				m.On("Save", mock.Anything).Return(nil)
			},
			requestBody: RequestSectionErrorJSON{
				SectionNumber:      123,
				MinimumTemperature: 15.0,
				CurrentCapacity:    50,
				MinimumCapacity:    30,
				MaximumCapacity:    100,
				WarehouseID:        1,
				ProductTypeID:      2,
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   *resterr.NewUnprocessableEntityError("current_temperature is required"),
		},
		{
			name:               "erro ao fazer decode",
			mockSetup:          func(m *MockSectionService) {},
			requestBody:        "invalid json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   *resterr.NewBadRequestError("json: cannot unmarshal string into Go value of type map[string]interface {}"),
		},
		{
			name: "return fail error when required field is missing",
			mockSetup: func(m *MockSectionService) {
				m.On("Save", mock.Anything).Return(nil)
			},
			requestBody: handler.RequestSectionJSON{
				SectionNumber:      123,
				CurrentTemperature: 22.5,
				MinimumTemperature: 15.0,
				CurrentCapacity:    3,
				MinimumCapacity:    30,
				MaximumCapacity:    100,
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   *resterr.NewUnprocessableEntityError("couldn't parse section"),
		},
		{
			name: "return conflict error when number is already in use",
			mockSetup: func(m *MockSectionService) {
				m.On("FindAll").Return([]internal.Section{
					{SectionNumber: 123},
				}, nil)
				m.On("Save", mock.Anything).Return(internal.ErrSectionNumberAlreadyInUse)
			},
			requestBody: handler.RequestSectionJSON{
				SectionNumber:      123,
				CurrentTemperature: 22.5,
				MinimumTemperature: 15.0,
				CurrentCapacity:    50,
				MinimumCapacity:    30,
				MaximumCapacity:    100,
				WarehouseID:        1,
				ProductTypeID:      2,
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   *resterr.NewConflictError("section with given section number already registered"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSectionService)
			sectionHandler := handler.NewHandlerSection(mockService)
			tt.mockSetup(mockService)

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/sections", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("content-type", "application/json")

			rr := httptest.NewRecorder()
			hd := sectionHandler.Create
			hd(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case internal.Section:
					var actualResponse = struct {
						Data internal.Section `json:"data"`
					}{}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response, actualResponse.Data)
				case resterr.RestErr:
					var actualResponse resterr.RestErr
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

func TestHandler_ReadAllSectionUnitTest(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSectionService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "successfully read all sections",
			mockSetup: func(m *MockSectionService) {
				sections := []internal.Section{
					{
						ID:                 1,
						SectionNumber:      123,
						CurrentTemperature: 22.5,
						MinimumTemperature: 15.0,
						CurrentCapacity:    50,
						MinimumCapacity:    30,
						MaximumCapacity:    100,
						WarehouseID:        1,
						ProductTypeID:      2,
					},
					{
						ID:                 2,
						SectionNumber:      456,
						CurrentTemperature: 22.5,
						MinimumTemperature: 15.0,
						CurrentCapacity:    50,
						MinimumCapacity:    30,
						MaximumCapacity:    100,
						WarehouseID:        2,
						ProductTypeID:      2,
					},
				}
				m.On("FindAll").Return(sections, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": []internal.Section{
					{
						ID:                 1,
						SectionNumber:      123,
						CurrentTemperature: 22.5,
						MinimumTemperature: 15.0,
						CurrentCapacity:    50,
						MinimumCapacity:    30,
						MaximumCapacity:    100,
						WarehouseID:        1,
						ProductTypeID:      2,
					},
					{
						ID:                 2,
						SectionNumber:      456,
						CurrentTemperature: 22.5,
						MinimumTemperature: 15.0,
						CurrentCapacity:    50,
						MinimumCapacity:    30,
						MaximumCapacity:    100,
						WarehouseID:        2,
						ProductTypeID:      2,
					},
				},
			},
		},
		{
			name: "return error when no sections found",
			mockSetup: func(m *MockSectionService) {
				m.On("FindAll").Return([]internal.Section{}, internal.ErrSectionNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   resterr.NewNotFoundError("section not found"),
		},
		{
			name: "return internal server error",
			mockSetup: func(m *MockSectionService) {
				m.On("FindAll").Return([]internal.Section{}, errors.New("internal server error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSectionService)
			sectionHandler := handler.NewHandlerSection(mockService)
			tt.mockSetup(mockService)

			req, err := http.NewRequest(http.MethodGet, "/sections", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			hd := sectionHandler.GetAll
			hd(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data []internal.Section `json:"data"`
					}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					require.ElementsMatch(t, response["data"], actualResponse.Data)

				case *resterr.RestErr:
					var actualResponse *resterr.RestErr
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					require.Equal(t, response, actualResponse)

				default:
					t.Fatalf("unexpected response type: %T", response)
				}
			}
		})
	}
}

func TestHandler_ReadByIdSectionUnitTest(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSectionService)
		id                 string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should return a section by ID",
			mockSetup: func(m *MockSectionService) {
				section := internal.Section{
					ID:                 1,
					SectionNumber:      123,
					CurrentTemperature: 22.5,
					MinimumTemperature: 15.0,
					CurrentCapacity:    50,
					MinimumCapacity:    30,
					MaximumCapacity:    100,
					WarehouseID:        1,
					ProductTypeID:      2,
				}
				m.On("FindByID", 1).Return(section, nil)
			},
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": internal.Section{
					ID:                 1,
					SectionNumber:      123,
					CurrentTemperature: 22.5,
					MinimumTemperature: 15.0,
					CurrentCapacity:    50,
					MinimumCapacity:    30,
					MaximumCapacity:    100,
					WarehouseID:        1,
					ProductTypeID:      2,
				},
			},
		},
		{
			name: "should return not found error",
			mockSetup: func(m *MockSectionService) {
				m.On("FindByID", 1).Return(internal.Section{}, internal.ErrSectionNotFound)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *resterr.NewNotFoundError("section not found"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockSectionService) {
				m.On("FindByID", 1).Return(internal.Section{}, errors.New("internal server error"))
			},
			id:                 "1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name:               "should return bad request error",
			mockSetup:          func(m *MockSectionService) {},
			id:                 "invalid_id",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSectionService)
			sectionHandler := handler.NewHandlerSection(mockService)
			tt.mockSetup(mockService)

			req, err := http.NewRequest(http.MethodGet, "/sections/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			hd := sectionHandler.GetByID
			hd(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data internal.Section `json:"data"`
					}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)
				case resterr.RestErr:
					var actualResponse resterr.RestErr
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

func TestHandler_ReportProductsUnitTest(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSectionService)
		queryID            string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should return report for all sections",
			mockSetup: func(m *MockSectionService) {
				m.On("ReportProducts").Return([]internal.ReportProduct{
					{SectionID: 1, SectionNumber: 123, ProductsCount: 10},
					{SectionID: 2, SectionNumber: 456, ProductsCount: 5},
				}, nil)
			},
			queryID:            "",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": []internal.ReportProduct{
					{SectionID: 1, SectionNumber: 123, ProductsCount: 10},
					{SectionID: 2, SectionNumber: 456, ProductsCount: 5},
				},
			},
		},
		{
			name: "should return report for specific section",
			mockSetup: func(m *MockSectionService) {
				m.On("ReportProductsByID", 1).Return(internal.ReportProduct{
					SectionID:     1,
					SectionNumber: 123,
					ProductsCount: 10,
				}, nil)
			},
			queryID:            "1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": []internal.ReportProduct{
					{
						SectionID:     1,
						SectionNumber: 123,
						ProductsCount: 10,
					},
				},
			},
		},
		{
			name:               "should return bad request error for invalid ID",
			mockSetup:          func(m *MockSectionService) {},
			queryID:            "invalid_id",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
		{
			name: "should return not found error when the section ID does not exist",
			mockSetup: func(m *MockSectionService) {
				m.On("ReportProductsByID", 2).Return(internal.ReportProduct{}, internal.ErrSectionNotFound)
			},
			queryID:            "2",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *resterr.NewNotFoundError("section not found"),
		},
		{
			name: "should return not found error when the section ID does not exist",
			mockSetup: func(m *MockSectionService) {
				m.On("ReportProductsByID", 2).Return(internal.ReportProduct{}, internal.ErrSectionUnprocessableEntity)
			},
			queryID:            "2",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   *resterr.NewInternalServerError("couldn't parse section"),
		},
		{
			name: "should return report product",
			mockSetup: func(m *MockSectionService) {
				m.On("ReportProducts").Return([]internal.ReportProduct{}, internal.ErrReportProductNotFound)
			},
			queryID:            "",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   *resterr.NewInternalServerError("report product not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSectionService)
			sectionHandler := handler.NewHandlerSection(mockService)
			tt.mockSetup(mockService)

			url := "/sections/report-products"

			if tt.queryID != "" {
				url += "?id=" + tt.queryID
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			hd := sectionHandler.ReportProducts
			hd(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data []internal.ReportProduct `json:"data"`
					}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)
				case resterr.RestErr:
					var actualResponse resterr.RestErr
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response, actualResponse)
				default:
					t.Fatalf("unexpected response type: %T", response)
				}
			}
		})
	}
}

func TestHandler_UpdateSectionUnitTest(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSectionService)
		id                 string
		requestBody        interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should update a section",
			mockSetup: func(m *MockSectionService) {
				mockSection := internal.Section{
					ID:                 1,
					SectionNumber:      123,
					CurrentTemperature: 22.5,
					MinimumTemperature: 15.0,
					CurrentCapacity:    50,
					MinimumCapacity:    30,
					MaximumCapacity:    100,
					WarehouseID:        1,
					ProductTypeID:      2,
				}
				m.On("Update", 1, internal.SectionPatch{
					SectionNumber:      intPtr(123),
					CurrentTemperature: float64Ptr(22.5),
					MinimumTemperature: float64Ptr(15.0),
					CurrentCapacity:    intPtr(50),
					MinimumCapacity:    intPtr(30),
					MaximumCapacity:    intPtr(100),
					WarehouseID:        intPtr(1),
					ProductTypeID:      intPtr(2),
				}).Return(mockSection, nil)
			},
			id: "1",
			requestBody: handler.SectionsUpdateJSON{
				SectionNumber:      intPtr(123),
				CurrentTemperature: float64Ptr(22.5),
				MinimumTemperature: float64Ptr(15.0),
				CurrentCapacity:    intPtr(50),
				MinimumCapacity:    intPtr(30),
				MaximumCapacity:    intPtr(100),
				WarehouseID:        intPtr(1),
				ProductTypeID:      intPtr(2),
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"data": internal.Section{
					ID:                 1,
					SectionNumber:      123,
					CurrentTemperature: 22.5,
					MinimumTemperature: 15.0,
					CurrentCapacity:    50,
					MinimumCapacity:    30,
					MaximumCapacity:    100,
					WarehouseID:        1,
					ProductTypeID:      2,
				},
			},
		},
		{
			name: "should return bad request error for invalid id",
			mockSetup: func(m *MockSectionService) {
			},
			id:                 "invalid_id",
			requestBody:        nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
		{
			name: "erro ao fazer decode",
			mockSetup: func(m *MockSectionService) {
			},
			id:                 "1",
			requestBody:        "invalid json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   *resterr.NewBadRequestError("json: cannot unmarshal string into Go value of type handler.SectionsUpdateJSON"),
		},
		{
			name: "should return unprocessable entity error for invalid request body",
			mockSetup: func(m *MockSectionService) {
				m.On("Update", 1, mock.Anything).Return(internal.Section{}, internal.ErrSectionUnprocessableEntity)
			},
			id: "1",
			requestBody: handler.SectionsUpdateJSON{
				SectionNumber:      intPtr(456),
				CurrentTemperature: float64Ptr(22.5),
				MinimumTemperature: float64Ptr(15.0),
				CurrentCapacity:    intPtr(100),
				MinimumCapacity:    intPtr(30),
				MaximumCapacity:    intPtr(100),
				WarehouseID:        intPtr(1),
				ProductTypeID:      intPtr(2),
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   *resterr.NewUnprocessableEntityError("couldn't parse section"),
		},
		{
			name: "should return conflict error",
			mockSetup: func(m *MockSectionService) {
				m.On("Update", 1, mock.Anything).Return(internal.Section{}, internal.ErrSectionNumberAlreadyInUse)
			},
			id: "1",
			requestBody: handler.SectionsUpdateJSON{
				SectionNumber:      intPtr(456),
				CurrentTemperature: float64Ptr(22.5),
				MinimumTemperature: float64Ptr(15.0),
				CurrentCapacity:    intPtr(150),
				MinimumCapacity:    intPtr(30),
				MaximumCapacity:    intPtr(100),
				WarehouseID:        intPtr(1),
				ProductTypeID:      intPtr(2),
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   *resterr.NewConflictError("section with given section number already registered"),
		},
		{
			name: "should return not found error",
			mockSetup: func(m *MockSectionService) {
				m.On("Update", 1, mock.Anything).Return(internal.Section{}, internal.ErrSectionNotFound)
			},
			id: "1",
			requestBody: handler.SectionsUpdateJSON{
				SectionNumber:      intPtr(456),
				CurrentTemperature: float64Ptr(22.5),
				MinimumTemperature: float64Ptr(15.0),
				CurrentCapacity:    intPtr(150),
				MinimumCapacity:    intPtr(30),
				MaximumCapacity:    intPtr(100),
				WarehouseID:        intPtr(1),
				ProductTypeID:      intPtr(2),
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *resterr.NewNotFoundError("section not found"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockSectionService) {
				m.On("Update", 1, mock.Anything).Return(internal.Section{}, errors.New("internal server error"))
			},
			id: "1",
			requestBody: handler.SectionsUpdateJSON{
				SectionNumber:      intPtr(456),
				CurrentTemperature: float64Ptr(22.5),
				MinimumTemperature: float64Ptr(15.0),
				CurrentCapacity:    intPtr(150),
				MinimumCapacity:    intPtr(30),
				MaximumCapacity:    intPtr(100),
				WarehouseID:        intPtr(1),
				ProductTypeID:      intPtr(2),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSectionService)
			sectionHandler := handler.NewHandlerSection(mockService)
			tt.mockSetup(mockService)

			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, "/sections/"+tt.id, bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("content-type", "application/json")

			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			hd := sectionHandler.Update
			hd(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case map[string]interface{}:
					var actualResponse struct {
						Data internal.Section `json:"data"`
					}
					err = json.NewDecoder(rr.Body).Decode(&actualResponse)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t, response["data"], actualResponse.Data)
				case resterr.RestErr:
					var actualResponse resterr.RestErr
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

func TestHandler_DeleteSectionUnitTest(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSectionService)
		id                 string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should delete a section",
			mockSetup: func(m *MockSectionService) {
				m.On("Delete", 1).Return(nil)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   nil,
		},
		{
			name: "should return not found error",
			mockSetup: func(m *MockSectionService) {
				m.On("Delete", 1).Return(internal.ErrSectionNotFound)
			},
			id:                 "1",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   *resterr.NewNotFoundError("section not found"),
		},
		{
			name: "should return internal server error",
			mockSetup: func(m *MockSectionService) {
				m.On("Delete", 1).Return(errors.New("internal server error"))
			},
			id:                 "1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   nil,
		},
		{
			name:               "should return bad request error",
			mockSetup:          func(m *MockSectionService) {},
			id:                 "invalid_id",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSectionService)
			sectionHandler := handler.NewHandlerSection(mockService)
			tt.mockSetup(mockService)

			req, err := http.NewRequest(http.MethodDelete, "/sections/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			hd := sectionHandler.Delete
			hd(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedResponse != nil {
				switch response := tt.expectedResponse.(type) {
				case resterr.RestErr:
					var actualResponse resterr.RestErr
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
