package service_test

import (
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func NewWarehouseRepositoryMock() *WarehouseRepositoryMock {
	return &WarehouseRepositoryMock{}
}

type WarehouseRepositoryMock struct {
	mock.Mock
}

func (r *WarehouseRepositoryMock) FindAll() ([]internal.Warehouse, error) {
	args := r.Called()
	return args.Get(0).([]internal.Warehouse), args.Error(1)
}

func (r *WarehouseRepositoryMock) FindByID(id int) (internal.Warehouse, error) {
	args := r.Called(id)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func (r *WarehouseRepositoryMock) Save(warehouse *internal.Warehouse) error {
	args := r.Called(warehouse)
	return args.Error(0)
}

func (r *WarehouseRepositoryMock) Update(warehouse *internal.Warehouse) error {
	args := r.Called(warehouse)
	return args.Error(0)
}

func (r *WarehouseRepositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func TestWarehouseServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WarehouseServiceTestSuite))
}

type WarehouseServiceTestSuite struct {
	rp *WarehouseRepositoryMock
	sv *service.WarehouseDefault
	suite.Suite
}

func (s *WarehouseServiceTestSuite) SetupTest() {
	rp := NewWarehouseRepositoryMock()
	sv := service.NewWarehouseDefault(rp)

	s.rp = rp
	s.sv = sv
}

// TestCases is a struct that contains the test cases for the warehouse service
func (w *WarehouseServiceTestSuite) TestWarehouseService_Save() {
	w.T().Run("case 1 - success: Should save a warehouse", func(t *testing.T) {
		w.SetupTest()
		warehouse := internal.Warehouse{
			ID:                 1,
			WarehouseCode:      "warehouse_code",
			Address:            "address",
			Telephone:          "phone",
			MinimumCapacity:    1,
			MinimumTemperature: 1,
		}
		w.rp.On("FindAll").Return([]internal.Warehouse{}, nil)
		w.rp.On("Save", &warehouse).Return(nil)

		err := w.sv.Save(&warehouse)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindAll", 1)
		w.rp.AssertNumberOfCalls(w.T(), "Save", 1)
		w.NoError(err)
	})

	w.T().Run("case 2 - error: Should return an error when trying to save a warehouse with duplicated code", func(t *testing.T) {
		w.SetupTest()
		warehouse := internal.Warehouse{
			ID:                 1,
			WarehouseCode:      "warehouse_code",
			Address:            "address",
			Telephone:          "phone",
			MinimumCapacity:    1,
			MinimumTemperature: 1,
		}
		w.rp.On("FindAll").Return([]internal.Warehouse{warehouse}, nil)

		err := w.sv.Save(&warehouse)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindAll", 1)
		w.rp.AssertNumberOfCalls(w.T(), "Save", 0)
		w.Equal(internal.ErrWarehouseRepositoryDuplicated, err)
	})

	w.T().Run("case 3 - error: Should return an error when trying to save a warehouse with invalid data", func(t *testing.T) {
		w.SetupTest()
		warehouse := internal.Warehouse{
			ID:                 1,
			WarehouseCode:      "warehouse_code",
			Address:            "address",
			Telephone:          "phone",
			MinimumCapacity:    0,
			MinimumTemperature: 1,
		}

		err := w.sv.Save(&warehouse)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindAll", 0)
		w.rp.AssertNumberOfCalls(w.T(), "Save", 0)
		w.Error(err, internal.ErrWarehouseUnprocessableEntity)
	})

	w.T().Run("case 4 - error: Should return an error when an internal error occurs", func(t *testing.T) {
		w.SetupTest()
		warehouse := internal.Warehouse{
			ID:                 1,
			WarehouseCode:      "warehouse_code",
			Address:            "address",
			Telephone:          "phone",
			MinimumCapacity:    1,
			MinimumTemperature: 1,
		}
		w.rp.On("FindAll").Return([]internal.Warehouse{}, errors.New("internal server error"))

		err := w.sv.Save(&warehouse)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindAll", 1)
		w.rp.AssertNumberOfCalls(w.T(), "Save", 0)
		w.Error(err)

	})

}

func (w *WarehouseServiceTestSuite) TestWarehouseService_FindAll() {
	warehouseSlice := []internal.Warehouse{
		{
			ID:                 1,
			WarehouseCode:      "warehouse_code",
			Address:            "address",
			Telephone:          "phone",
			MinimumCapacity:    1,
			MinimumTemperature: 1,
		},
		{
			ID:                 2,
			WarehouseCode:      "warehouse_code",
			Address:            "address",
			Telephone:          "phone",
			MinimumCapacity:    1,
			MinimumTemperature: 1,
		},
	}
	w.rp.On("FindAll").Return(warehouseSlice, nil)

	warehouses, err := w.sv.FindAll()
	require.NoError(w.T(), err)
	require.Equal(w.T(), warehouseSlice, warehouses)
}

func (w *WarehouseServiceTestSuite) TestWarehouseService_FindByID() {
	w.T().Run("case 1 - success: Should return a warehouse by id", func(t *testing.T) {
		w.SetupTest()
		warehouse := internal.Warehouse{
			ID:                 1,
			WarehouseCode:      "warehouse_code",
			Address:            "address",
			Telephone:          "phone",
			MinimumCapacity:    1,
			MinimumTemperature: 1,
		}
		w.rp.On("FindByID", 1).Return(warehouse, nil)

		result, err := w.sv.FindByID(1)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindByID", 1)
		require.NoError(w.T(), err)
		require.Equal(w.T(), warehouse, result)
	})

	w.T().Run("case 2 - error: Should return an error when trying to find a warehouse by a non existent id", func(t *testing.T) {
		w.SetupTest()
		w.rp.On("FindByID", 2).Return(internal.Warehouse{}, internal.ErrWarehouseRepositoryNotFound)

		_, err := w.sv.FindByID(2)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindByID", 1)
		require.Error(w.T(), err)
		w.Equal(internal.ErrWarehouseRepositoryNotFound, err)
	})
}

func (w *WarehouseServiceTestSuite) TestWarehouseService_Update() {
	w.T().Run("case 1 - success: Should update a warehouse", func(t *testing.T) {
		w.SetupTest()
		// attributes
		warehouseCode := "W1"
		address := "123 Main St"
		telephone := "123-456-7890"
		minimumCapacity := 100
		minimumTemperature := 20.5

		// warehouse patch
		warehousePatch := internal.WarehousePatchUpdate{
			WarehouseCode:      &warehouseCode,
			Address:            &address,
			Telephone:          &telephone,
			MinimumCapacity:    &minimumCapacity,
			MinimumTemperature: &minimumTemperature,
		}
		// warehouse for update
		warehouse := internal.Warehouse{
			ID:                 1,
			WarehouseCode:      "warehouse_code",
			Address:            "address",
			Telephone:          "phone",
			MinimumCapacity:    1,
			MinimumTemperature: 1,
		}
		// updated warehouse
		warehouseUpdated := internal.Warehouse{
			ID:                 1,
			WarehouseCode:      "W1",
			Address:            "123 Main St",
			Telephone:          "123-456-7890",
			MinimumCapacity:    100,
			MinimumTemperature: 20.5,
		}

		w.rp.On("FindByID", 1).Return(warehouse, nil)
		w.rp.On("FindAll").Return([]internal.Warehouse{}, nil)
		w.rp.On("Update", mock.AnythingOfType("*internal.Warehouse")).Run(func(args mock.Arguments) {
			w := args.Get(0).(*internal.Warehouse)
			w.WarehouseCode = "W1"
			w.Address = "123 Main St"
			w.Telephone = "123-456-7890"
			w.MinimumCapacity = 100
			w.MinimumTemperature = 20.5
		}).Return(nil)

		result, err := w.sv.Update(1, &warehousePatch)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindByID", 1)
		w.rp.AssertNumberOfCalls(w.T(), "Update", 1)
		require.NoError(w.T(), err)
		require.Equal(w.T(), warehouseUpdated, result)
	})

	w.T().Run("case 2 - error: Should return an error when trying to update a warehouse with a duplicated code", func(t *testing.T) {
		w.SetupTest()
		// attributes
		warehouseCode := "W1"
		address := "123 Main St"
		telephone := "123-456-7890"
		minimumCapacity := 100
		minimumTemperature := 20.5

		// warehouse patch
		warehousePatch := internal.WarehousePatchUpdate{
			WarehouseCode:      &warehouseCode,
			Address:            &address,
			Telephone:          &telephone,
			MinimumCapacity:    &minimumCapacity,
			MinimumTemperature: &minimumTemperature,
		}
		// warehouse with duplicated code
		warehouseDuplicated := internal.Warehouse{
			ID:                 2,
			WarehouseCode:      "W1",
			Address:            "123 Main St",
			Telephone:          "123-456-7890",
			MinimumCapacity:    100,
			MinimumTemperature: 20.5,
		}

		w.rp.On("FindByID", 1).Return(internal.Warehouse{}, nil)
		w.rp.On("FindAll").Return([]internal.Warehouse{warehouseDuplicated}, nil)

		_, err := w.sv.Update(1, &warehousePatch)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindByID", 1)
		w.rp.AssertNumberOfCalls(w.T(), "FindAll", 1)
		w.rp.AssertNumberOfCalls(w.T(), "Update", 0)
		w.Equal(internal.ErrWarehouseRepositoryDuplicated, err)
	})

	w.T().Run("case 3 - error: Should return an error when trying to update a non existent warehouse", func(t *testing.T) {
		w.SetupTest()
		warehousePatch := internal.WarehousePatchUpdate{}
		w.rp.On("FindByID", 2).Return(internal.Warehouse{}, internal.ErrWarehouseRepositoryNotFound)

		_, err := w.sv.Update(2, &warehousePatch)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindByID", 1)
		w.rp.AssertNumberOfCalls(w.T(), "Update", 0)
		w.Equal(internal.ErrWarehouseRepositoryNotFound, err)
	})

	w.T().Run("case 4 - error: Should return an error when an internal error occurs", func(t *testing.T) {
		w.SetupTest()
		// attributes
		warehouseCode := "W1"
		address := "123 Main St"
		telephone := "123-456-7890"
		minimumCapacity := 100
		minimumTemperature := 20.5

		// warehouse patch
		warehousePatch := internal.WarehousePatchUpdate{
			WarehouseCode:      &warehouseCode,
			Address:            &address,
			Telephone:          &telephone,
			MinimumCapacity:    &minimumCapacity,
			MinimumTemperature: &minimumTemperature,
		}
		// warehouse for update
		warehouse := internal.Warehouse{
			ID:                 1,
			WarehouseCode:      "warehouse_code",
			Address:            "address",
			Telephone:          "phone",
			MinimumCapacity:    1,
			MinimumTemperature: 1,
		}

		w.rp.On("FindByID", 1).Return(warehouse, nil)
		w.rp.On("FindAll").Return([]internal.Warehouse{}, errors.New("internal server error"))

		_, err := w.sv.Update(1, &warehousePatch)
		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindByID", 1)
		w.rp.AssertNumberOfCalls(w.T(), "FindAll", 1)
		w.rp.AssertNumberOfCalls(w.T(), "Update", 0)
		require.Error(w.T(), err)
	})
}

func (w *WarehouseServiceTestSuite) TestWarehouseService_Delete() {
	w.T().Run("case 1 - success: Should delete a warehouse", func(t *testing.T) {
		w.SetupTest()

		w.rp.On("FindByID", 1).Return(internal.Warehouse{}, nil)
		w.rp.On("Delete", 1).Return(nil)

		err := w.sv.Delete(1)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindByID", 1)
		w.rp.AssertNumberOfCalls(w.T(), "Delete", 1)
		require.NoError(w.T(), err)
	})

	w.T().Run("case 2 - error: Should return an error when trying to delete a non existent warehouse", func(t *testing.T) {
		w.SetupTest()

		w.rp.On("FindByID", 2).Return(internal.Warehouse{}, internal.ErrWarehouseRepositoryNotFound)

		err := w.sv.Delete(2)

		w.rp.AssertExpectations(w.T())
		w.rp.AssertNumberOfCalls(w.T(), "FindByID", 1)
		w.rp.AssertNumberOfCalls(w.T(), "Delete", 0)
		require.Error(w.T(), err)
		w.Equal(internal.ErrWarehouseRepositoryNotFound, err)
	})
}
