package repository_test

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestRepository_NewMapEmployeeUnitTest(t *testing.T) {
	t.Run("creates map successfully", func(t *testing.T) {
		mp, e := repository.NewEmployeeRepository("../../db/employees.json")

		require.NoError(t, e)
		require.NotNil(t, mp)
	})
	t.Run("fails to create map, invalid path", func(t *testing.T) {
		mp, e := repository.NewEmployeeRepository("not a valid path")

		require.Nil(t, mp)
		require.Error(t, e)
	})
	t.Run("fails to create map, invalid json structure", func(t *testing.T) {
		mp, e := repository.NewEmployeeRepository("../../db/employees_test.json")

		require.Nil(t, mp)
		require.Error(t, e)
	})
}

func TestRepository_MapImplementationsEmployeeUnitTest(t *testing.T) {
	mp, e := repository.NewEmployeeRepository("../../db/employees.json")

	require.NoError(t, e)
	require.NotNil(t, mp)
	expectedEmployees := map[int]internal.Employee{
		1: {
			ID:           1,
			CardNumberID: "E1001",
			FirstName:    "John",
			LastName:     "Doe",
			WarehouseID:  1,
		},
		2: {
			ID:           2,
			CardNumberID: "E1002",
			FirstName:    "Jane",
			LastName:     "Smith",
			WarehouseID:  2,
		},
		3: {
			ID:           3,
			CardNumberID: "E1003",
			FirstName:    "Michael",
			LastName:     "Johnson",
			WarehouseID:  3,
		},
		4: {
			ID:           4,
			CardNumberID: "E1004",
			FirstName:    "Emily",
			LastName:     "Davis",
			WarehouseID:  4,
		},
		5: {
			ID:           5,
			CardNumberID: "E1005",
			FirstName:    "David",
			LastName:     "Miller",
			WarehouseID:  5,
		},
	}

	t.Run("getAll", func(t *testing.T) {
		actualEmployees := mp.GetAll()

		require.Equal(t, expectedEmployees, actualEmployees)
	})
	t.Run("add", func(t *testing.T) {
		employee := internal.Employee{
			ID:           0,
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			CardNumberID: "3456",
			WarehouseID:  6,
		}
		id := mp.Save(&employee)
		rp := mp.GetAll()

		require.Equal(t, employee, rp[id])
	})
	t.Run("update", func(t *testing.T) {
		expectedEmployee := internal.Employee{
			ID:           5,
			CardNumberID: "1234",
			FirstName:    "404",
			LastName:     "NoName",
		}
		mp.Update(5, expectedEmployee)
		rp := mp.GetAll()

		require.Equal(t, expectedEmployee, rp[5])
	})
	t.Run("delete", func(t *testing.T) {
		mp.Delete(5)
		rp := mp.GetAll()

		_, exists := rp[5]
		require.False(t, exists)
	})
}
