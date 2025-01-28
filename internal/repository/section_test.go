package repository_test

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestRepository_NewMapSectionUnitTest(t *testing.T) {
	t.Run("creates map successfully", func(t *testing.T) {
		mp, e := repository.NewRepositorySection("../../db/section.json")

		require.NoError(t, e)
		require.NotNil(t, mp)
	})
	t.Run("fails to create map, invalid path", func(t *testing.T) {
		mp, e := repository.NewRepositorySection("not a valid path")

		require.Nil(t, mp)
		require.Error(t, e)
	})
	t.Run("fails to create map, invalid json structure", func(t *testing.T) {
		mp, e := repository.NewRepositorySection("../../db/section_test.json")

		require.Nil(t, mp)
		require.Error(t, e)
	})
}

func TestRepository_MapImplementationsSectionUnitTest(t *testing.T) {
	mp, e := repository.NewRepositorySection("../../db/section.json")

	require.NoError(t, e)
	require.NotNil(t, mp)
	expectedSections := []internal.Section{
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
		{
			ID:                 2,
			SectionNumber:      102,
			CurrentTemperature: 18.0,
			MinimumTemperature: 10.0,
			CurrentCapacity:    40,
			MinimumCapacity:    25,
			MaximumCapacity:    85,
			WarehouseID:        202,
			ProductTypeID:      302,
		},
	}

	t.Run("getAll", func(t *testing.T) {
		actualSections, e := mp.FindAll()

		require.NoError(t, e)
		require.Equal(t, expectedSections, actualSections)
	})
	t.Run("getAll (empty)", func(t *testing.T) {
		mp, e := repository.NewRepositorySection("../../db/empty_section_test.json")

		require.NoError(t, e)
		require.NotNil(t, mp)

		actualSections, e := mp.FindAll()

		require.Error(t, e)
		require.Nil(t, actualSections)
	})
	t.Run("findByID success", func(t *testing.T) {
		section, e := mp.FindByID(1)

		require.NoError(t, e)
		require.Equal(t, expectedSections[0], section)
	})
	t.Run("findByID failure", func(t *testing.T) {
		section, e := mp.FindByID(3)

		require.Error(t, e)
		require.Zero(t, section)
	})
	t.Run("save success", func(t *testing.T) {
		expectedSection := internal.Section{
			SectionNumber:      102,
			CurrentTemperature: 18.0,
			MinimumTemperature: 10.0,
			CurrentCapacity:    40,
			MinimumCapacity:    25,
			MaximumCapacity:    85,
			WarehouseID:        202,
			ProductTypeID:      302,
		}
		e := mp.Save(&expectedSection)

		require.NoError(t, e)

		actualSection, e := mp.FindByID(expectedSection.ID)

		require.NoError(t, e)
		require.Equal(t, expectedSection, actualSection)
	})
	t.Run("update success", func(t *testing.T) {
		expectedSection := internal.Section{
			ID:                 3,
			SectionNumber:      106,
			CurrentTemperature: 19.0,
			MinimumTemperature: 10.0,
			CurrentCapacity:    40,
			MinimumCapacity:    25,
			MaximumCapacity:    85,
			WarehouseID:        202,
			ProductTypeID:      302,
		}
		e := mp.Update(&expectedSection)

		require.NoError(t, e)

		actualSection, e := mp.FindByID(expectedSection.ID)

		require.NoError(t, e)
		require.Equal(t, expectedSection, actualSection)
	})
	t.Run("update failure", func(t *testing.T) {
		expectedSection := internal.Section{
			ID: 10,
		}
		e := mp.Update(&expectedSection)

		require.Error(t, e)
	})
	t.Run("delete success", func(t *testing.T) {
		e := mp.Delete(3)

		require.NoError(t, e)
	})
	t.Run("delete failure", func(t *testing.T) {
		e := mp.Delete(10)

		require.Error(t, e)
	})
	t.Run("check if section number exists (it doesnt)", func(t *testing.T) {
		e := mp.SectionNumberExists(internal.Section{
			ID:            3,
			SectionNumber: 103,
		})

		require.NoError(t, e)
	})
	t.Run("check if section number exists (it does)", func(t *testing.T) {
		e := mp.SectionNumberExists(internal.Section{
			ID:            3,
			SectionNumber: 101,
		})

		require.Error(t, e)
	})
	t.Run("reportProducts (unimplemented)", func(t *testing.T) {
		mp.ReportProducts()
	})
	t.Run("reportProductsByID (unimplemented)", func(t *testing.T) {
		mp.ReportProductsByID(-1)
	})
}
