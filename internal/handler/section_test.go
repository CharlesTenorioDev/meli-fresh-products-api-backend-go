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

func ok(s internal.Section) bool {
	if s.SectionNumber <= 0 ||
		s.CurrentTemperature < -273.15 ||
		s.MinimumTemperature < -273.15 ||
		s.CurrentCapacity < 0 ||
		s.MinimumCapacity < 0 ||
		s.MaximumCapacity < 0 ||
		s.WarehouseID <= 0 ||
		s.ProductTypeID <= 0 {
		return false
	}

	return true
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
	args := m.Called()
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

func (m *MockSectionService) Update(id int, updates map[string]interface{}) (internal.Section, error) {
	args := m.Called(id, updates)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (m *MockSectionService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetAllSections(t *testing.T) {
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

func TestGetSectionById(t *testing.T) {
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

func TestSaveSection(t *testing.T) {
	tests := []struct {
		name               string
		mockSetup          func(*MockSectionService)
		requestBody        interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "should create a new section",
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
			name: "should return conflict error when number is already in use",
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
		{
			name: "return unprocessable entity error when warehouse ID does not exist",
			mockSetup: func(m *MockSectionService) {
				m.On("Save", mock.Anything).Return(nil)
			},
			requestBody: handler.RequestSectionJSON{
				SectionNumber:      123,
				CurrentTemperature: 22.5,
				MinimumTemperature: 15.0,
				CurrentCapacity:    -3,
				MinimumCapacity:    30,
				MaximumCapacity:    100,
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse:   *resterr.NewUnprocessableEntityError("couldn't parse section"),
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

func TestUpdateSection(t *testing.T) {

}

func TestDeleteSection(t *testing.T) {
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

/*
func TestUpdateSection() {
	section := internal.Section{
		SectionNumber:      101,
		CurrentTemperature: 1.0,
		MinimumTemperature: 1.0,
		CurrentCapacity:    1,
		MinimumCapacity:    1,
		MaximumCapacity:    1,
		WarehouseID:        1,
		ProductTypeID:      1,
	}

	suite.service.On("Update", 1, mock.Anything).Return(section, nil)
	suite.service.On("FindByID", 1).Return(section, nil)

	body, _ := json.Marshal(section)
	r := httptest.NewRequest(http.MethodPut, "/sections/{id}", bytes.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.Update(w, r)

	assert.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)

	var response struct {
		Data internal.Section `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), section, response.Data)
}
*/
