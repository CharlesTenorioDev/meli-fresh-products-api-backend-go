package service_test

import (
	"fmt"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type employeeRepositoryMock struct {
	mock.Mock
}

func (r *employeeRepositoryMock) GetAll() (db []internal.Employee, err error) {
	args := r.Called()
	return args.Get(0).([]internal.Employee), args.Error(1)
}

func (r *employeeRepositoryMock) GetByID(id int) (emp internal.Employee, err error) {
	args := r.Called(id)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (r *employeeRepositoryMock) Save(emp *internal.Employee) (int, error) {
	args := r.Called(emp)
	return args.Get(0).(int), args.Error(1)
}

func (r *employeeRepositoryMock) Update(id int, employee internal.Employee) error {
	args := r.Called(id, employee)
	return args.Error(0)
}

func (r *employeeRepositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *employeeRepositoryMock) CountInboundOrdersPerEmployee() (io []internal.InboundOrdersPerEmployee, err error) {
	args := r.Called()
	return args.Get(0).([]internal.InboundOrdersPerEmployee), args.Error(1)
}

func (r *employeeRepositoryMock) ReportInboundOrdersByID(employeeID int) (io internal.InboundOrdersPerEmployee, err error) {
	args := r.Called(employeeID)
	return args.Get(0).(internal.InboundOrdersPerEmployee), args.Error(1)
}

func TestCreate_Employee(t *testing.T) {
	defaultEmployee := internal.Employee{
		ID:           0,
		CardNumberID: "abcdef",
		FirstName:    "Fabio",
		LastName:     "Nacarelli",
		WarehouseID:  14,
	}
	t.Run("create an employee successfully", func(t *testing.T) {
		rpWarehouse := new(warehouseRepositoryMock)
		rp := new(employeeRepositoryMock)
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rpWarehouse.On("FindByID", 14).Return(internal.Warehouse{
			ID: 14,
		}, nil)
		rp.On("GetAll").Return([]internal.Employee{}, nil)
		rp.On("Save", &defaultEmployee).Return(defaultEmployee.ID, nil)
		err := sv.Save(&defaultEmployee)

		require.NoError(t, err)
	})
	t.Run("create an employee with conflict", func(t *testing.T) {
		rpWarehouse := new(warehouseRepositoryMock)
		rp := new(employeeRepositoryMock)
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return([]internal.Employee{
			{
				CardNumberID: "abcdef",
			},
		}, nil)
		err := sv.Save(&defaultEmployee)

		require.ErrorIs(t, err, service.ErrCardNumberIDInUse)
	})
}

func TestRead_Employee(t *testing.T) {
	expectedEmployees := []internal.Employee{
		{
			ID:        1,
			FirstName: "Fabio",
			LastName:  "Nacarelli",
		},
		{
			ID:        2,
			FirstName: "Mocked",
			LastName:  "Database",
		},
	}
	t.Run("read every employee", func(t *testing.T) {
		rpWarehouse := new(warehouseRepositoryMock)
		rp := new(employeeRepositoryMock)
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return(expectedEmployees, nil)

		employees, err := sv.GetAll()

		require.NoError(t, err)
		require.Equal(t, expectedEmployees, employees)
	})
	t.Run("read employee with id 1 (does not exist)", func(t *testing.T) {
		rpWarehouse := new(warehouseRepositoryMock)
		rp := new(employeeRepositoryMock)
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetByID", 1).Return(internal.Employee{}, fmt.Errorf("employee not found"))

		_, err := sv.GetByID(1)

		require.Error(t, err)
	})
	t.Run("read employee with id 2 (does exist)", func(t *testing.T) {
		expectedEmployee := internal.Employee{
			ID:        2,
			FirstName: "Mocked",
			LastName:  "Database",
		}
		rpWarehouse := new(warehouseRepositoryMock)
		rp := new(employeeRepositoryMock)
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetByID", 2).Return(expectedEmployee, nil)

		employee, err := sv.GetByID(2)

		require.NoError(t, err)
		require.Equal(t, expectedEmployee, employee)
	})
}

func TestUpdate_Employee(t *testing.T) {
	t.Run("update employee with id 1 (does not exist)", func(t *testing.T) {
		rpWarehouse := new(warehouseRepositoryMock)
		rp := new(employeeRepositoryMock)
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return([]internal.Employee{}, nil)
		employee := internal.Employee{
			ID:       1,
			LastName: "Fabio",
		}

		err := sv.Update(employee)
		require.Error(t, err)
	})
	t.Run("update employee with id 1 (does exist)", func(t *testing.T) {
		employee := internal.Employee{
			ID:           1,
			FirstName:    "Another",
			LastName:     "Name",
			CardNumberID: "fedcba",
			WarehouseID:  14,
		}
		rpWarehouse := new(warehouseRepositoryMock)
		rp := new(employeeRepositoryMock)
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rpWarehouse.On("FindByID", 14).Return(internal.Warehouse{
			ID: 14,
		}, nil)
		rp.On("GetAll").Return([]internal.Employee{
			{
				ID:           1,
				FirstName:    "Fabio",
				LastName:     "Nacarelli",
				CardNumberID: "abcdef",
				WarehouseID:  14,
			},
		}, nil)
		rp.On("Update", 1, employee).Return(nil)

		err := sv.Update(employee)

		require.NoError(t, err)
	})
}

func TestDelete_Employee(t *testing.T) {
	t.Run("user does not exist", func(t *testing.T) {
		rpWarehouse := new(warehouseRepositoryMock)
		rp := new(employeeRepositoryMock)
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetByID", 1).Return(internal.Employee{}, internal.ErrEmployeeNotFound)

		err := sv.Delete(1)

		require.ErrorIs(t, err, service.ErrEmployeeNotFound)
	})
	t.Run("user does exist", func(t *testing.T) {
		rpWarehouse := new(warehouseRepositoryMock)
		rp := new(employeeRepositoryMock)
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetByID", 1).Return(internal.Employee{}, nil)
		rp.On("Delete", 1).Return(nil)

		err := sv.Delete(1)
		require.NoError(t, err)
	})
}