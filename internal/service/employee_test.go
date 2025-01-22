package service_test

import (
	"fmt"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func NewEmployeeRepositoryMock() *EmployeeRepositoryMock {
	return &EmployeeRepositoryMock{}
}

type EmployeeRepositoryMock struct {
	mock.Mock
}

func (r *EmployeeRepositoryMock) GetAll() (db []internal.Employee, err error) {
	args := r.Called()
	return args.Get(0).([]internal.Employee), args.Error(1)
}

func (r *EmployeeRepositoryMock) GetByID(id int) (emp internal.Employee, err error) {
	args := r.Called(id)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (r *EmployeeRepositoryMock) Save(emp *internal.Employee) (int, error) {
	args := r.Called(emp)
	return args.Get(0).(int), args.Error(1)
}

func (r *EmployeeRepositoryMock) Update(id int, employee internal.Employee) error {
	args := r.Called(id, employee)
	return args.Error(0)
}

func (r *EmployeeRepositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *EmployeeRepositoryMock) CountInboundOrdersPerEmployee() (io []internal.InboundOrdersPerEmployee, err error) {
	args := r.Called()
	return args.Get(0).([]internal.InboundOrdersPerEmployee), args.Error(1)
}

func (r *EmployeeRepositoryMock) ReportInboundOrdersByID(employeeID int) (io internal.InboundOrdersPerEmployee, err error) {
	args := r.Called(employeeID)
	return args.Get(0).(internal.InboundOrdersPerEmployee), args.Error(1)
}

func TestCreate_EmployeeUnitTestService(t *testing.T) {
	defaultEmployee := internal.Employee{
		ID:           0,
		CardNumberID: "abcdef",
		FirstName:    "Fabio",
		LastName:     "Nacarelli",
		WarehouseID:  14,
	}
	t.Run("create an employee successfully", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
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
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return([]internal.Employee{
			{
				CardNumberID: "abcdef",
			},
		}, nil)
		err := sv.Save(&defaultEmployee)

		require.ErrorIs(t, err, service.ErrCardNumberIDInUse)
	})
	t.Run("create an employee with warehouse that does not exist", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rpWarehouse.On("FindByID", 14).Return(internal.Warehouse{}, internal.ErrWarehouseRepositoryNotFound)
		rp.On("GetAll").Return([]internal.Employee{}, nil)
		err := sv.Save(&defaultEmployee)

		require.Error(t, err)
	})
	t.Run("create an employee without the required fields", func(t *testing.T) {
		emp := internal.Employee{
			ID:           0,
			CardNumberID: "",
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			WarehouseID:  14,
		}
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return([]internal.Employee{}, nil)
		err := sv.Save(&emp)

		require.Error(t, err)
	})
}

func TestRead_EmployeeUnitTestService(t *testing.T) {
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
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return(expectedEmployees, nil)

		employees, err := sv.GetAll()

		require.NoError(t, err)
		require.Equal(t, expectedEmployees, employees)
	})
	t.Run("read employee with id 1 (does not exist)", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
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
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetByID", 2).Return(expectedEmployee, nil)

		employee, err := sv.GetByID(2)

		require.NoError(t, err)
		require.Equal(t, expectedEmployee, employee)
	})
	t.Run("reading every employee fails", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return([]internal.Employee{}, internal.ErrEmployeeNotFound)

		_, err := sv.GetAll()

		require.Error(t, err)
	})
}

func TestUpdate_EmployeeUnitTestService(t *testing.T) {
	t.Run("update employee with id 1 (does not exist)", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
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
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
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

func TestDelete_EmployeeUnitTestService(t *testing.T) {
	t.Run("user does not exist", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetByID", 1).Return(internal.Employee{}, internal.ErrEmployeeNotFound)

		err := sv.Delete(1)

		require.ErrorIs(t, err, service.ErrEmployeeNotFound)
	})
	t.Run("user does exist", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetByID", 1).Return(internal.Employee{}, nil)
		rp.On("Delete", 1).Return(nil)

		err := sv.Delete(1)
		require.NoError(t, err)
	})
}
