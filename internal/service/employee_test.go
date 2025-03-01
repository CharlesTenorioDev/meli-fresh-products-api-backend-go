package service_test

import (
	"errors"
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

func (r *EmployeeRepositoryMock) Save(emp *internal.Employee) (id int64, err error) {
	args := r.Called(emp)
	return args.Get(0).(int64), args.Error(1)
}

func (r *EmployeeRepositoryMock) Update(id int, employee internal.Employee) (err error) {
	args := r.Called(id, employee)
	return args.Error(0)
}

func (r *EmployeeRepositoryMock) Delete(id int) (err error) {
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

func TestCreate_EmployeeUnitTest(t *testing.T) {
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
		rp.On("Save", &defaultEmployee).Return(int64(defaultEmployee.ID), nil)
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
	t.Run("internal fails to get all employees", func(t *testing.T) {
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
		rp.On("GetAll").Return([]internal.Employee{}, internal.ErrEmployeeNotFound)
		err := sv.Save(&emp)

		require.Error(t, err)
	})
	t.Run("employee already exists", func(t *testing.T) {
		emp := internal.Employee{
			ID:           1,
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
	t.Run("create an employee fails", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rpWarehouse.On("FindByID", 14).Return(internal.Warehouse{
			ID: 14,
		}, nil)
		rp.On("GetAll").Return([]internal.Employee{}, nil)
		rp.On("Save", &defaultEmployee).Return(int64(-1), errors.New("failed to create employee"))
		err := sv.Save(&defaultEmployee)

		require.Error(t, err)
	})
}

func TestRead_EmployeeUnitTest(t *testing.T) {
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

func TestUpdate_EmployeeUnitTest(t *testing.T) {
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
	t.Run("update but fails to fetch employees", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return([]internal.Employee{}, internal.ErrEmployeeNotFound)
		employee := internal.Employee{
			ID:       1,
			LastName: "Fabio",
		}

		err := sv.Update(employee)
		require.Error(t, err)
	})
	t.Run("card number id already in use", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return([]internal.Employee{
			{
				ID:           1,
				CardNumberID: "abcdef",
			},
			{
				ID:           2,
				CardNumberID: "cdef",
			},
		}, nil)
		employee := internal.Employee{
			ID:           2,
			CardNumberID: "abcdef",
		}

		err := sv.Update(employee)
		require.Error(t, err)
	})
	t.Run("missing required fields", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetAll").Return([]internal.Employee{
			{
				ID:           1,
				CardNumberID: "abcdef",
			},
		}, nil)
		employee := internal.Employee{
			ID:           1,
			CardNumberID: "",
		}

		err := sv.Update(employee)
		require.Error(t, err)
	})
	t.Run("warehouse conflict", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rpWarehouse.On("FindByID", 1).Return(internal.Warehouse{}, service.ErrConflictInEmployee)
		rp.On("GetAll").Return([]internal.Employee{
			{
				ID:           1,
				CardNumberID: "abcdef",
			},
		}, nil)
		employee := internal.Employee{
			ID:           1,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			WarehouseID:  1,
			CardNumberID: "abcd",
		}

		err := sv.Update(employee)
		require.Error(t, err)
	})
	t.Run("updating fails", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		employee := internal.Employee{
			ID:           1,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			WarehouseID:  1,
			CardNumberID: "abcd",
		}
		rpWarehouse.On("FindByID", 1).Return(internal.Warehouse{}, nil)
		rp.On("GetAll").Return([]internal.Employee{
			{
				ID:           1,
				CardNumberID: "abcdef",
			},
		}, nil)
		rp.On("Update", 1, employee).Return(errors.New("internal err update"))

		err := sv.Update(employee)
		require.Error(t, err)
	})
}

func TestDelete_EmployeeUnitTest(t *testing.T) {
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
	t.Run("user does not exist pt2", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("GetByID", 1).Return(internal.Employee{}, errors.New("just coverage"))

		err := sv.Delete(1)

		require.Error(t, err)
	})
}

func Test_CoverageEmployeeUnitTest(t *testing.T) {
	t.Run("call count inbound orders per employee", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("CountInboundOrdersPerEmployee").Return([]internal.InboundOrdersPerEmployee{}, nil)

		sv.CountInboundOrdersPerEmployee()
	})
	t.Run("call report inbound orders per employee", func(t *testing.T) {
		rpWarehouse := NewWarehouseRepositoryMock()
		rp := NewEmployeeRepositoryMock()
		sv := service.NewEmployeeServiceDefault(rp, rpWarehouse)
		rp.On("ReportInboundOrdersByID", 1).Return(internal.InboundOrdersPerEmployee{}, nil)

		sv.ReportInboundOrdersByID(1)
	})
}
