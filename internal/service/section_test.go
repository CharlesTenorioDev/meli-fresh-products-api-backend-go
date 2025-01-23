package service_test

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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

func intPtr(i int) *int {
	return &i
}

func newSectionService() (*service.SectionService, *SectionRepositoryMock, *service.ProductTypeRepositoryMock, *WarehouseRepositoryMock) {
	rpSection := new(SectionRepositoryMock)
	rpProductType := new(service.ProductTypeRepositoryMock)
	rpProduct := new(repositoryProductMock)
	rpWareHouse := new(WarehouseRepositoryMock)

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
		rpProductType.AssertExpectations(t)
		rpWareHouse.AssertExpectations(t)
	})

	t.Run("return fail error when required field is missing", func(t *testing.T) {
		sv, rpSection, _, rpWareHouse := newSectionService()

		sectionCreate := newTestSection(0, 101, 0, 2)

		rpSection.On("SectionNumberExists", sectionCreate.SectionNumber).Return(false, nil)
		rpWareHouse.On("FindByID", sectionCreate.WarehouseID).Return(internal.Warehouse{}, internal.ErrWarehouseRepositoryNotFound)

		err := sv.Save(&sectionCreate)

		require.ErrorIs(t, err, internal.ErrSectionUnprocessableEntity)
		require.Contains(t, err.Error(), "couldn't parse section")
		//rpSection.AssertExpectations(t)
		//rpWareHouse.AssertExpectations(t)
	})

	t.Run("return conflict error when number is already in use", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		sectionCreate := newTestSection(0, 101, 4, 3)

		rpSection.On("SectionNumberExists", sectionCreate.SectionNumber).Return(true, nil)

		err := sv.Save(&sectionCreate)

		require.ErrorIs(t, err, internal.ErrSectionNumberAlreadyInUse)
		rpSection.AssertExpectations(t)
		rpProductType.AssertExpectations(t)
		rpWareHouse.AssertExpectations(t)
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
	})

	t.Run("return error when reading a nonexistent section by ID", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()
		expectedError := internal.ErrSectionNotFound

		rpSection.On("FindByID", 1).Return(internal.Section{}, expectedError)

		_, err := sv.FindByID(1)

		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)
		rpSection.AssertExpectations(t)
	})

	t.Run("successfully read an existing section by ID", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()
		expectedSection := newTestSection(2, 101, 4, 3)

		rpSection.On("FindByID", 2).Return(expectedSection, nil)

		section, err := sv.FindByID(2)

		require.NoError(t, err)
		require.Equal(t, expectedSection, section)
		rpSection.AssertExpectations(t)
	})
}

func TestService_UpdateSectionUnitTest(t *testing.T) {
	t.Run("return error when updating a nonexistent section", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()

		updates := internal.SectionPatch{
			CurrentCapacity: intPtr(50),
		}

		rpSection.On("FindByID", 1).Return(internal.Section{}, internal.ErrSectionNotFound)

		updatedSection, err := sv.Update(1, updates)

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionNotFound, err)
		require.Empty(t, updatedSection)
		rpSection.AssertExpectations(t)
	})

	t.Run("successfully update an existing section", func(t *testing.T) {
		sv, rpSection, rpProductType, rpWareHouse := newSectionService()

		existingSection := newTestSection(1, 100, 4, 3)

		updates := internal.SectionPatch{
			CurrentCapacity: intPtr(150),
		}

		rpSection.On("FindByID", 1).Return(existingSection, nil)
		rpWareHouse.On("FindByID", existingSection.WarehouseID).Return(internal.Warehouse{ID: existingSection.WarehouseID}, nil)
		rpProductType.On("FindByID", existingSection.ProductTypeID).Return(internal.ProductType{ID: existingSection.ProductTypeID}, nil)

		rpSection.On("Update", mock.AnythingOfType("*internal.Section")).Return(nil)

		updatedSection, err := sv.Update(1, updates)

		require.NoError(t, err)
		require.NotEqual(t, existingSection, updatedSection)
		rpSection.AssertExpectations(t)
	})
}

func TestService_DeleteSectionUnitTest(t *testing.T) {
	t.Run("return error when attempting to delete a nonexistent section", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()

		rpSection.On("FindByID", 1).Return(internal.Section{}, internal.ErrSectionNotFound)

		err := sv.Delete(1)

		require.ErrorIs(t, err, internal.ErrSectionNotFound)
		rpSection.AssertExpectations(t)
	})

	t.Run("successfully delete an existing section", func(t *testing.T) {
		sv, rpSection, _, _ := newSectionService()

		existingSection := newTestSection(1, 101, 4, 3)

		rpSection.On("FindByID", 1).Return(existingSection, nil)
		rpSection.On("Delete", 1).Return(nil)

		err := sv.Delete(1)

		require.NoError(t, err)
		rpSection.AssertExpectations(t)
	})
}
