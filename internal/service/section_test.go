package service_test

import (
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func NewSectionRepositoryMock() *SectionRepositoryMock {
	return &SectionRepositoryMock{}
}

type SectionRepositoryMock struct {
	mock.Mock
}

func (r *SectionRepositoryMock) FindAll() ([]internal.Section, error) {
	args := r.Called()
	return args.Get(0).([]internal.Section), args.Error(1)
}

func (r *SectionRepositoryMock) FindByID(id int) (internal.Section, error) {
	args := r.Called(id)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (r *SectionRepositoryMock) ReportProducts() ([]internal.ReportProduct, error) {
	args := r.Called()
	return args.Get(0).([]internal.ReportProduct), args.Error(1)
}

func (r *SectionRepositoryMock) ReportProductsByID(sectionID int) (internal.ReportProduct, error) {
	args := r.Called(sectionID)
	return args.Get(0).(internal.ReportProduct), args.Error(1)
}

func (r *SectionRepositoryMock) SectionNumberExists(sectionNumber int) (bool, error) {
	args := r.Called(sectionNumber)
	return args.Get(0).(bool), args.Error(1)
}

func (r *SectionRepositoryMock) Save(section *internal.Section) error {
	if section.WarehouseID == 0 {
		return internal.ErrWarehouseRepositoryNotFound
	}

	args := r.Called(section)
	return args.Error(0)
}

func (r *SectionRepositoryMock) Update(section *internal.Section) error {
	args := r.Called(section)
	return args.Error(0)
}

func (r *SectionRepositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func float64Ptr(f float64) *float64 {
	return &f
}

func newSectionService() (*service.SectionService, *SectionRepositoryMock, *ProductTypeRepositoryMock, *WarehouseRepositoryMock) {
	rpSection := NewSectionRepositoryMock()
	rpProductType := NewProductTypeRepositoryMock()
	rpProduct := NewRepositoryProductMock()
	rpWareHouse := NewWarehouseRepositoryMock()

	return service.NewServiceSection(rpSection, rpProductType, rpProduct, rpWareHouse), rpSection, rpProductType, rpWareHouse
}

func newTestSection(id int, sectionNumber int, warehouseID int, productTypeID int) internal.Section {
	return internal.Section{
		ID:                 id,
		SectionNumber:      sectionNumber,
		CurrentTemperature: 22.5,
		MinimumTemperature: 15.0,
		CurrentCapacity:    50,
		MinimumCapacity:    30,
		MaximumCapacity:    100,
		WarehouseID:        warehouseID,
		ProductTypeID:      productTypeID,
	}
}

func TestService_CreateSectionUnitTest(t *testing.T) {
	t.Run("successfully create a new section", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		sectionCreate := newTestSection(0, 101, 4, 3)

		rpSection.On("SectionNumberExists", sectionCreate.SectionNumber).Return(false, nil)
		rpWareHouse.On("FindByID", sectionCreate.WarehouseID).Return(internal.Warehouse{ID: sectionCreate.WarehouseID}, nil)
		rpProductType.On("FindByID", sectionCreate.ProductTypeID).Return(internal.ProductType{ID: sectionCreate.ProductTypeID}, nil)
		rpSection.On("Save", &sectionCreate).Return(nil)

		err := sv.Save(&sectionCreate)

		require.NoError(t, err)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "Save", 1)
	})

	t.Run("return fail error when required field is missing", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		sectionCreate := newTestSection(0, 101, 0, 2)

		err := sv.Save(&sectionCreate)

		require.ErrorIs(t, err, internal.ErrSectionUnprocessableEntity)
		require.Contains(t, err.Error(), "couldn't parse section")

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 0)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Save", 0)

	})

	t.Run("return conflict error when number is already in use", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		sectionCreate := newTestSection(0, 101, 4, 3)

		rpSection.On("SectionNumberExists", sectionCreate.SectionNumber).Return(true, nil)

		err := sv.Save(&sectionCreate)

		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrSectionNumberAlreadyInUse)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Save", 0)
	})

	t.Run("returns error when warehouse does not exist", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		sectionCreate := newTestSection(0, 101, 99, 3)

		rpSection.On("SectionNumberExists", sectionCreate.SectionNumber).Return(false, nil)
		rpWareHouse.On("FindByID", sectionCreate.WarehouseID).Return(internal.Warehouse{}, internal.ErrWarehouseRepositoryNotFound)

		err := sv.Save(&sectionCreate)

		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrWarehouseRepositoryNotFound)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Save", 0)
	})

	t.Run("returns error when product type does not exist", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		sectionCreate := newTestSection(0, 101, 3, 99)

		rpSection.On("SectionNumberExists", sectionCreate.SectionNumber).Return(false, nil)
		rpWareHouse.On("FindByID", sectionCreate.WarehouseID).Return(internal.Warehouse{}, nil)
		rpProductType.On("FindByID", sectionCreate.ProductTypeID).Return(internal.ProductType{}, internal.ErrProductTypeNotFound)

		err := sv.Save(&sectionCreate)

		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrProductTypeNotFound)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "Save", 0)
	})

	t.Run("returns error when section fails to save", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		sectionCreate := newTestSection(0, 101, 3, 99)

		rpSection.On("SectionNumberExists", sectionCreate.SectionNumber).Return(false, nil)
		rpWareHouse.On("FindByID", sectionCreate.WarehouseID).Return(internal.Warehouse{}, nil)
		rpProductType.On("FindByID", sectionCreate.ProductTypeID).Return(internal.ProductType{}, nil)
		rpSection.On("Save", &sectionCreate).Return(internal.ErrSectionUnprocessableEntity)

		err := sv.Save(&sectionCreate)

		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrSectionUnprocessableEntity)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "Save", 1)
	})
}

func TestService_ReadSectionUnitTest(t *testing.T) {
	t.Run("successfully read all sections", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()

		sectionsRead := []internal.Section{
			newTestSection(0, 101, 2, 1),
			newTestSection(0, 102, 1, 4),
		}

		rpSection.On("FindAll").Return(sectionsRead, nil)

		sections, err := sv.FindAll()

		require.NoError(t, err)
		require.Equal(t, sectionsRead, sections)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindAll", 1)
	})

	t.Run("returns an error when no section exists", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()
		expectedError := internal.ErrSectionNotFound

		rpSection.On("FindAll").Return([]internal.Section{}, expectedError)

		_, err := sv.FindAll()

		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindAll", 1)
	})

	t.Run("return error when reading a nonexistent section by ID", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()
		expectedError := internal.ErrSectionNotFound

		rpSection.On("FindByID", 1).Return(internal.Section{}, expectedError)

		_, err := sv.FindByID(1)

		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
	})

	t.Run("successfully read an existing section by ID", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()
		expectedSection := newTestSection(2, 101, 4, 3)

		rpSection.On("FindByID", 2).Return(expectedSection, nil)

		section, err := sv.FindByID(2)

		require.NoError(t, err)
		require.Equal(t, expectedSection, section)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
	})
}

func TestService_ReportProductsSectionUnitTest(t *testing.T) {
	t.Run("successfully report products", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()

		expectedReports := []internal.ReportProduct{
			{
				SectionID:     1,
				SectionNumber: 123,
				ProductsCount: 3,
			},
			{
				SectionID:     2,
				SectionNumber: 456,
				ProductsCount: 2,
			},
			{
				SectionID:     2,
				SectionNumber: 789,
				ProductsCount: 4,
			},
		}

		rpSection.On("ReportProducts").Return(expectedReports, nil)

		reports, err := sv.ReportProducts()

		require.NoError(t, err)
		require.Equal(t, expectedReports, reports)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "ReportProducts", 1)
	})

	t.Run("returns an error when no report products exists", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()
		expectedError := internal.ErrReportProductNotFound

		rpSection.On("ReportProducts").Return([]internal.ReportProduct{}, expectedError)

		_, err := sv.ReportProducts()

		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "ReportProducts", 1)
	})
}

func TestService_ReportProductsByIDSectionUnitTest(t *testing.T) {
	t.Run("successfully report products by section ID", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()
		sectionID := 1
		expectedReport := internal.ReportProduct{
			SectionID:     1,
			SectionNumber: 123,
			ProductsCount: 3,
		}

		rpSection.On("FindByID", sectionID).Return(internal.Section{}, nil)
		rpSection.On("ReportProductsByID", sectionID).Return(expectedReport, nil)

		report, err := sv.ReportProductsByID(sectionID)

		require.NoError(t, err)
		require.Equal(t, expectedReport, report)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "ReportProductsByID", 1)
	})

	t.Run("return error when section ID does not exist", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()
		sectionID := 1
		expectedError := internal.ErrSectionNotFound

		rpSection.On("FindByID", sectionID).Return(internal.Section{}, expectedError)

		report, err := sv.ReportProductsByID(sectionID)

		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)
		require.Empty(t, report)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "ReportProductsByID", 0)
	})

	t.Run("return error when reporting products by section ID fails", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()
		sectionID := 1
		expectedError := errors.New("error on reporting products")

		rpSection.On("FindByID", sectionID).Return(internal.Section{}, nil)
		rpSection.On("ReportProductsByID", sectionID).Return(internal.ReportProduct{}, expectedError)

		report, err := sv.ReportProductsByID(sectionID)

		require.Error(t, err)
		require.EqualError(t, err, expectedError.Error())
		require.Empty(t, report)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "ReportProductsByID", 1)
	})
}

func TestService_UpdateSectionUnitTest(t *testing.T) {
	t.Run("return error when updating a nonexistent section", func(t *testing.T) {
		sv, rpSection, rpWareHouse, rpProductType := newSectionService()

		updates := internal.SectionPatch{
			CurrentCapacity: intPtr(150),
		}

		rpSection.On("FindByID", 1).Return(internal.Section{}, internal.ErrSectionNotFound)

		updatedSection, err := sv.Update(1, updates)

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionNotFound, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 0)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("successfully update an existing section", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		existingSection := newTestSection(1, 100, 6, 7)

		updates := internal.SectionPatch{
			SectionNumber:      intPtr(456),
			CurrentTemperature: float64Ptr(22.5),
			MinimumTemperature: float64Ptr(15.0),
			CurrentCapacity:    intPtr(502),
			MinimumCapacity:    intPtr(302),
			MaximumCapacity:    intPtr(150),
			WarehouseID:        intPtr(7),
			ProductTypeID:      intPtr(7),
		}

		rpSection.On("FindByID", 1).Return(existingSection, nil)
		rpSection.On("SectionNumberExists", *updates.SectionNumber).Return(false, nil)
		rpWareHouse.On("FindByID", *updates.WarehouseID).Return(internal.Warehouse{ID: *updates.WarehouseID}, nil)
		rpProductType.On("FindByID", *updates.ProductTypeID).Return(internal.ProductType{ID: *updates.ProductTypeID}, nil)

		rpSection.On("Update", mock.AnythingOfType("*internal.Section")).Return(nil)

		updatedSection, err := sv.Update(1, updates)

		require.NoError(t, err)
		require.NotEqual(t, existingSection, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "Update", 1)
	})

	t.Run("returns error when find section number", func(t *testing.T) {
		sv, rpSection, rpWareHouse, rpProductType := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, nil)
		rpSection.On("SectionNumberExists", 123).Return(false, errors.New("error when find section number"))

		updatedSection, err := sv.Update(1, internal.SectionPatch{SectionNumber: intPtr(123)})

		require.Error(t, err)
		require.Equal(t, errors.New("error when find section number"), err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when section number is already in use", func(t *testing.T) {
		sv, rpSection, rpWareHouse, rpProductType := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, nil)
		rpSection.On("SectionNumberExists", 123).Return(true, nil)

		updatedSection, err := sv.Update(1, internal.SectionPatch{SectionNumber: intPtr(123)})

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionNumberAlreadyInUse, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when section number is less than or equal to zero", func(t *testing.T) {
		sv, rpSection, rpWareHouse, rpProductType := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, nil)
		rpSection.On("SectionNumberExists", 0).Return(false, nil)

		updatedSection, err := sv.Update(1, internal.SectionPatch{SectionNumber: intPtr(0)})

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionUnprocessableEntity, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when current temperature is below absolute zero", func(t *testing.T) {
		sv, rpSection, rpWareHouse, rpProductType := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, nil)
		rpSection.On("SectionNumberExists", 123).Return(false, nil)

		updatedSection, err := sv.Update(1, internal.SectionPatch{SectionNumber: intPtr(123), CurrentTemperature: float64Ptr(-274)})

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionUnprocessableEntity, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when minimum temperature is below absolute zero", func(t *testing.T) {
		sv, rpSection, rpWareHouse, rpProductType := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, nil)
		rpSection.On("SectionNumberExists", 123).Return(false, nil)

		updatedSection, err := sv.Update(1, internal.SectionPatch{SectionNumber: intPtr(123), MinimumTemperature: float64Ptr(-274)})

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionUnprocessableEntity, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when current capacity is less than zero", func(t *testing.T) {
		sv, rpSection, rpWareHouse, rpProductType := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, nil)
		rpSection.On("SectionNumberExists", 123).Return(false, nil)

		updatedSection, err := sv.Update(1, internal.SectionPatch{SectionNumber: intPtr(123), CurrentCapacity: intPtr(-2)})

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionUnprocessableEntity, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when minimum capacity is less than zero", func(t *testing.T) {
		sv, rpSection, rpWareHouse, rpProductType := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, nil)
		rpSection.On("SectionNumberExists", 123).Return(false, nil)

		updatedSection, err := sv.Update(1, internal.SectionPatch{SectionNumber: intPtr(123), MinimumCapacity: intPtr(-2)})

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionUnprocessableEntity, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when maximum capacity is less than zero", func(t *testing.T) {
		sv, rpSection, rpWareHouse, rpProductType := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, nil)
		rpSection.On("SectionNumberExists", 123).Return(false, nil)

		updatedSection, err := sv.Update(1, internal.SectionPatch{SectionNumber: intPtr(123), MaximumCapacity: intPtr(-2)})

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionUnprocessableEntity, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when warehouse does not exist", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		existingSection := newTestSection(1, 100, 6, 7)

		updates := internal.SectionPatch{
			SectionNumber:      intPtr(456),
			CurrentTemperature: float64Ptr(22.5),
			MinimumTemperature: float64Ptr(15.0),
			CurrentCapacity:    intPtr(502),
			MinimumCapacity:    intPtr(302),
			MaximumCapacity:    intPtr(150),
			WarehouseID:        intPtr(7),
			ProductTypeID:      intPtr(7),
		}

		rpSection.On("FindByID", 1).Return(existingSection, nil)
		rpSection.On("SectionNumberExists", *updates.SectionNumber).Return(false, nil)
		rpWareHouse.On("FindByID", *updates.WarehouseID).Return(internal.Warehouse{ID: *updates.WarehouseID}, internal.ErrWarehouseRepositoryNotFound)

		updatedSection, err := sv.Update(1, updates)

		require.Error(t, err)
		require.Equal(t, internal.ErrWarehouseRepositoryNotFound, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when product type does not exist", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		existingSection := newTestSection(1, 100, 6, 7)

		updates := internal.SectionPatch{
			SectionNumber:      intPtr(456),
			CurrentTemperature: float64Ptr(22.5),
			MinimumTemperature: float64Ptr(15.0),
			CurrentCapacity:    intPtr(502),
			MinimumCapacity:    intPtr(302),
			MaximumCapacity:    intPtr(150),
			WarehouseID:        intPtr(7),
			ProductTypeID:      intPtr(7),
		}

		rpSection.On("FindByID", 1).Return(existingSection, nil)
		rpSection.On("SectionNumberExists", *updates.SectionNumber).Return(false, nil)
		rpWareHouse.On("FindByID", *updates.WarehouseID).Return(internal.Warehouse{ID: *updates.WarehouseID}, nil)
		rpProductType.On("FindByID", *updates.ProductTypeID).Return(internal.ProductType{ID: *updates.ProductTypeID}, internal.ErrProductTypeNotFound)

		updatedSection, err := sv.Update(1, updates)

		require.Error(t, err)
		require.Equal(t, internal.ErrProductTypeNotFound, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "Update", 0)
	})

	t.Run("returns error when section fails to update", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		existingSection := newTestSection(1, 100, 6, 7)

		updates := internal.SectionPatch{
			SectionNumber:      intPtr(456),
			CurrentTemperature: float64Ptr(22.5),
			MinimumTemperature: float64Ptr(15.0),
			CurrentCapacity:    intPtr(502),
			MinimumCapacity:    intPtr(302),
			MaximumCapacity:    intPtr(150),
			WarehouseID:        intPtr(7),
			ProductTypeID:      intPtr(7),
		}

		rpSection.On("FindByID", 1).Return(existingSection, nil)
		rpSection.On("SectionNumberExists", *updates.SectionNumber).Return(false, nil)
		rpWareHouse.On("FindByID", *updates.WarehouseID).Return(internal.Warehouse{ID: *updates.WarehouseID}, nil)
		rpProductType.On("FindByID", *updates.ProductTypeID).Return(internal.ProductType{ID: *updates.ProductTypeID}, nil)

		rpSection.On("Update", mock.AnythingOfType("*internal.Section")).Return(internal.ErrSectionUnprocessableEntity)

		updatedSection, err := sv.Update(1, updates)

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionUnprocessableEntity, err)
		require.Empty(t, updatedSection)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "SectionNumberExists", 1)
		rpWareHouse.AssertExpectations(t)
		rpWareHouse.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductType.AssertExpectations(t)
		rpProductType.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "Update", 1)
	})
}

func TestService_DeleteSectionUnitTest(t *testing.T) {
	t.Run("return error when attempting to delete a nonexistent section", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, internal.ErrSectionNotFound)

		err := sv.Delete(1)

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionNotFound, err)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "Delete", 0)
	})

	t.Run("successfully delete an existing section", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()

		existingSection := newTestSection(1, 101, 4, 3)

		rpSection.On("FindByID", 1).Return(existingSection, nil)
		rpSection.On("Delete", 1).Return(nil)

		err := sv.Delete(1)

		require.NoError(t, err)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "Delete", 1)
	})

	t.Run("return error when attempting to delete a nonexistent section", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, nil)
		rpSection.On("Delete", 1).Return(errors.New("error delete section"))

		err := sv.Delete(1)

		require.Error(t, err)
		require.Equal(t, errors.New("error delete section"), err)

		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertNumberOfCalls(t, "Delete", 1)
	})
}
