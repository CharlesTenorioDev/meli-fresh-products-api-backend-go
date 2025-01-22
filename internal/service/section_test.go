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

func (r *SectionRepositoryMock) SectionNumberExists(section internal.Section) (bool, error) {
	args := r.Called(section)
	return args.Get(0).(bool), args.Error(1)
}

func (r *SectionRepositoryMock) Save(section *internal.Section) error {
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

func newSectionService() (*service.SectionService, *SectionRepositoryMock) {
	rpSection := new(SectionRepositoryMock)
	rpProductType := new(service.ProductTypeRepositoryMock)
	rpProduct := new(repositoryProductMock)
	rpWareHouse := new(WarehouseRepositoryMock)

	return service.NewServiceSection(rpSection, rpProductType, rpProduct, rpWareHouse), rpSection
}

func TestCreate_Section(t *testing.T) {
	sectionCreate := internal.Section{
		ID:                 0,
		SectionNumber:      101,
		CurrentTemperature: 22.5,
		MinimumTemperature: 15.0,
		CurrentCapacity:    50,
		MinimumCapacity:    30,
		MaximumCapacity:    100,
		WarehouseID:        4,
		ProductTypeID:      3,
	}

	t.Run("successfully create a new section", func(t *testing.T) {
		sv, rpSection := newSectionService()

		rpSection.On("FindAll").Return([]internal.Section{}, nil)
		rpSection.On("Save", &sectionCreate).Return(1, nil)

		err := sv.Save(&sectionCreate)

		require.NoError(t, err)
		rpSection.AssertExpectations(t)
	})

	t.Run("conflict to create a section when number is already in use", func(t *testing.T) {
		sv, rpSection := newSectionService()

		rpSection.On("FindAll").Return([]internal.Section{
			{SectionNumber: 101},
		}, nil)

		err := sv.Save(&sectionCreate)

		require.ErrorIs(t, err, internal.ErrSectionNumberAlreadyInUse)
		rpSection.AssertExpectations(t)
	})
}

func TestRead_Section(t *testing.T) {
	sectionsRead := []internal.Section{
		{
			ID:                 0,
			SectionNumber:      101,
			CurrentTemperature: 22.5,
			MinimumTemperature: 15.0,
			CurrentCapacity:    50,
			MinimumCapacity:    30,
			MaximumCapacity:    100,
			WarehouseID:        2,
			ProductTypeID:      1,
		},
		{
			ID:                 0,
			SectionNumber:      102,
			CurrentTemperature: 29.5,
			MinimumTemperature: 19.0,
			CurrentCapacity:    56,
			MinimumCapacity:    36,
			MaximumCapacity:    102,
			WarehouseID:        1,
			ProductTypeID:      4,
		},
	}

	t.Run("successfully read all sections", func(t *testing.T) {
		sv, rpSection := newSectionService()
		rpSection.On("FindAll").Return(sectionsRead, nil)

		sections, err := sv.FindAll()

		require.NoError(t, err)
		require.Equal(t, sectionsRead, sections)
		rpSection.AssertExpectations(t)
	})

	t.Run("return error when reading a nonexistent section by ID", func(t *testing.T) {
		sv, rpSection := newSectionService()
		rpSection.On("FindByID", 1).Return(internal.Section{}, internal.ErrSectionNotFound)

		_, err := sv.FindByID(1)

		require.Error(t, err)
		require.EqualError(t, err, "section not found")
		rpSection.AssertExpectations(t)
	})

	t.Run("successfully read an existing section by ID", func(t *testing.T) {
		expectedSection := internal.Section{
			ID:                 0,
			SectionNumber:      101,
			CurrentTemperature: 22.5,
			MinimumTemperature: 15.0,
			CurrentCapacity:    50,
			MinimumCapacity:    30,
			MaximumCapacity:    100,
			WarehouseID:        4,
			ProductTypeID:      3,
		}

		sv, rpSection := newSectionService()
		rpSection.On("FindByID", 2).Return(expectedSection, nil)

		section, err := sv.FindByID(2)

		require.NoError(t, err)
		require.Equal(t, expectedSection, section)
		rpSection.AssertExpectations(t)
	})
}

func TestUpdate_Section(t *testing.T) {
	t.Run("return error when updating a nonexistent section", func(t *testing.T) {
		sv, rpSection := newSectionService()

		updates := map[string]interface{}{
			"maximum_capacity": 150,
			"warehouse_id":     5,
			"product_type_id":  4,
		}

		rpSection.On("Update", 1, updates).Return(internal.Section{}, internal.ErrSectionNotFound)

		updatedSection, err := sv.Update(1, updates)

		require.Error(t, err)
		require.Equal(t, internal.ErrSectionNotFound, err)
		require.Empty(t, updatedSection)
		rpSection.AssertExpectations(t)
	})

	t.Run("successfully update an existing section", func(t *testing.T) {
		sv, rpSection := newSectionService()

		existingSection := internal.Section{
			ID:              1,
			MaximumCapacity: 100,
			WarehouseID:     14,
			ProductTypeID:   2,
		}

		updates := map[string]interface{}{
			"maximum_capacity": 150,
		}

		rpSection.On("Update", 1, updates).Return(existingSection, nil)

		updatedSection, err := sv.Update(1, updates)

		require.NoError(t, err)
		require.Equal(t, existingSection, updatedSection)
		rpSection.AssertExpectations(t)
	})
}

func TestDelete_Section(t *testing.T) {
	t.Run("return error when attempting to delete a nonexistent section", func(t *testing.T) {
		sv, rpSection := newSectionService()
		rpSection.On("FindByID", 1).Return(internal.Section{}, internal.ErrSectionNotFound)

		err := sv.Delete(1)

		require.ErrorIs(t, err, internal.ErrSectionNotFound)
		rpSection.AssertExpectations(t)
	})
	t.Run("successfully delete an existing section", func(t *testing.T) {
		sv, rpSection := newSectionService()
		rpSection.On("FindByID", 1).Return(internal.Section{ID: 1}, nil)
		rpSection.On("Delete", 1).Return(nil)

		err := sv.Delete(1)

		require.NoError(t, err)
		rpSection.AssertExpectations(t)
	})
}
