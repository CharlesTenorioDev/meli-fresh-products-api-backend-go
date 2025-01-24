package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MockSectionService struct {
	mock.Mock
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
	args := m.Called(section)
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

type SectionTestSuite struct {
	suite.Suite
	handler *handler.SectionHandler
	service *MockSectionService
}

func (suite *SectionTestSuite) SetupTest() {
	suite.service = new(MockSectionService)
	suite.handler = handler.NewHandlerSection(suite.service)
}

func (suite *SectionTestSuite) TestGetAllSections() {
	sections := []internal.Section{
		{
			ID:                 1,
			SectionNumber:      101,
			CurrentTemperature: 22.5,
			MinimumTemperature: 15.0,
			CurrentCapacity:    50,
			MinimumCapacity:    30,
			MaximumCapacity:    100,
			WarehouseID:        201,
			ProductTypeID:      301,
		},
	}

	var err error
	suite.service.On("FindAll").Return(sections, err)

	r := httptest.NewRequest(http.MethodGet, "/sections", nil)
	w := httptest.NewRecorder()
	suite.handler.GetAll(w, r)

	assert.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)

	var response struct {
		Data []internal.Section `json:"data"`
	}
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)

	expected := []internal.Section{
		{
			ID:                 1,
			SectionNumber:      101,
			CurrentTemperature: 22.5,
			MinimumTemperature: 15.0,
			CurrentCapacity:    50,
			MinimumCapacity:    30,
			MaximumCapacity:    100,
			WarehouseID:        201,
			ProductTypeID:      301,
		},
	}
	assert.Equal(suite.T(), expected, response.Data)
}

func (suite *SectionTestSuite) TestGetSectionById() {
	section := internal.Section{
		ID:                 1,
		SectionNumber:      101,
		CurrentTemperature: 22.5,
		MinimumTemperature: 15.0,
		CurrentCapacity:    50,
		MinimumCapacity:    30,
		MaximumCapacity:    100,
		WarehouseID:        201,
		ProductTypeID:      301,
	}

	var err error
	suite.service.On("FindByID", 1).Return(section, err)

	r := httptest.NewRequest(http.MethodGet, "/sections/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.GetByID(w, r)

	assert.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)

	var response struct {
		Data internal.Section `json:"data"`
	}
	err = json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), section, response.Data)
}

func (suite *SectionTestSuite) TestGetSectionByIdNotFound() {
	suite.service.On("FindByID", 1).Return(internal.Section{}, internal.ErrSectionNotFound)

	r := httptest.NewRequest(http.MethodGet, "/sections/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.GetByID(w, r)

	assert.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)

	var response struct {
		Error string `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), "section not found", response.Error)
}

func (suite *SectionTestSuite) TestSaveSection() {
	section := internal.Section{
		ID:                 1,
		SectionNumber:      101,
		CurrentTemperature: 22.5,
		MinimumTemperature: 15.0,
		CurrentCapacity:    50,
		MinimumCapacity:    30,
		MaximumCapacity:    100,
		WarehouseID:        201,
		ProductTypeID:      301,
	}
	suite.service.On("Save", mock.AnythingOfType("*internal.Section")).Run(func(args mock.Arguments) {
		w := args.Get(0).(*internal.Section)
		w.ID = section.ID
	}).Return(nil)

	body, _ := json.Marshal(section)
	r := httptest.NewRequest(http.MethodPost, "/sections", bytes.NewReader(body))
	w := httptest.NewRecorder()

	suite.handler.Create(w, r)

	assert.Equal(suite.T(), http.StatusCreated, w.Result().StatusCode)

	var response struct {
		Data internal.Section
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), section, response.Data)
}

func (suite *SectionTestSuite) TestSaveSectionError() {

	section := internal.Section{
		SectionNumber:      111,
		CurrentTemperature: 1.0,
		MinimumTemperature: 1.0,
		CurrentCapacity:    1,
		MinimumCapacity:    1,
		MaximumCapacity:    1,
		WarehouseID:        1,
		ProductTypeID:      1,
	}
	suite.service.On("Save", &section).Return(internal.ErrSectionNotFound)

	body, _ := json.Marshal(section)
	r := httptest.NewRequest(http.MethodPost, "/sections", bytes.NewReader(body))
	w := httptest.NewRecorder()

	suite.handler.Create(w, r)

	assert.Equal(suite.T(), http.StatusUnprocessableEntity, w.Result().StatusCode)

	var response struct {
		Error string `json:"error"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response.Error)
	assert.Equal(suite.T(), "unprocessable_entity", response.Error)
}

func (suite *SectionTestSuite) TestUpdateSection() {
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

func (suite *SectionTestSuite) TestDeleteSection() {
	suite.service.On("Delete", 1).Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/sections/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.Delete(w, r)

	assert.Equal(suite.T(), http.StatusNoContent, w.Result().StatusCode)
}

func (suite *SectionTestSuite) TestDeleteSectionNotFound() {
	suite.service.On("Delete", 1).Return(internal.ErrSectionNotFound)

	r := httptest.NewRequest(http.MethodDelete, "/sections/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	suite.handler.Delete(w, r)

	assert.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)
	var response struct {
		Error string `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "", response.Error)
}

func TestSectionTestSuite(t *testing.T) {
	suite.Run(t, new(SectionTestSuite))
}
