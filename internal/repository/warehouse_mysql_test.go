package repository_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestWarehouse_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := `
		SELECT
			id, warehouse_code, address, telephone, minimum_capacity, minimum_temperature
		FROM
			warehouses;
	`

	t.Run("case 1: success - Warehouses found", func(t *testing.T) {
		expectedWs := []internal.Warehouse{
			{
				ID:                 1,
				WarehouseCode:      "123ABC",
				Address:            "address",
				Telephone:          "telephone",
				MinimumCapacity:    1,
				MinimumTemperature: 1,
			},
			{
				ID:                 2,
				WarehouseCode:      "456DEF",
				Address:            "address",
				Telephone:          "telephone",
				MinimumCapacity:    1,
				MinimumTemperature: 1,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimum_capacity", "minimum_temperature"}).
			AddRow(expectedWs[0].ID, expectedWs[0].WarehouseCode, expectedWs[0].Address, expectedWs[0].Telephone, expectedWs[0].MinimumCapacity, expectedWs[0].MinimumTemperature).
			AddRow(expectedWs[1].ID, expectedWs[1].WarehouseCode, expectedWs[1].Address, expectedWs[1].Telephone, expectedWs[1].MinimumCapacity, expectedWs[1].MinimumTemperature)

		mock.ExpectQuery(query).WillReturnRows(rows)

		rp := repository.NewWarehouseMysqlRepository(db)
		ws, err := rp.FindAll()

		require.NoError(t, err)
		require.Equal(t, expectedWs, ws)
	})

	t.Run("case 2: error - Error executing the query", func(t *testing.T) {
		mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)

		rp := repository.NewWarehouseMysqlRepository(db)
		_, err := rp.FindAll()

		require.Error(t, err)
	})

	t.Run("case 3: error - Error iterating over the rows", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "123ABC", "address", "telephone", 1, 1).
			RowError(0, sql.ErrConnDone)

		mock.ExpectQuery(query).WillReturnRows(rows)

		rp := repository.NewWarehouseMysqlRepository(db)
		_, err := rp.FindAll()

		require.Error(t, err)
	})
}

func TestWarehouse_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := `
		SELECT
			id, warehouse_code, address, telephone, minimum_capacity, minimum_temperature
		FROM
			warehouses
		WHERE
			id = ?
	`

	t.Run("case 1: success - Warehouse found", func(t *testing.T) {
		id := 1
		expectedW := internal.Warehouse{
			ID:                 1,
			WarehouseCode:      "123ABC",
			Address:            "address",
			Telephone:          "telephone",
			MinimumCapacity:    1,
			MinimumTemperature: 1,
		}

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimum_capacity", "minimum_temperature"}).
			AddRow(expectedW.ID, expectedW.WarehouseCode, expectedW.Address, expectedW.Telephone, expectedW.MinimumCapacity, expectedW.MinimumTemperature)

		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(rows)

		rp := repository.NewWarehouseMysqlRepository(db)
		w, err := rp.FindByID(id)

		require.NoError(t, err)
		require.Equal(t, expectedW, w)
	})

	t.Run("case 2: error - Warehouse not found", func(t *testing.T) {
		id := 1
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		rp := repository.NewWarehouseMysqlRepository(db)
		_, err := rp.FindByID(id)

		require.Error(t, err)
		require.Equal(t, internal.ErrWarehouseRepositoryNotFound, err)
	})

	t.Run("case 3: error - Error executing the query", func(t *testing.T) {
		id := 1
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnError(sql.ErrConnDone)

		rp := repository.NewWarehouseMysqlRepository(db)
		_, err := rp.FindByID(id)

		require.Error(t, err)
	})
}

func TestWarehouseMysql_Save(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	query := `
		INSERT INTO warehouses (warehouse_code, address, telephone, minimum_capacity, minimum_temperature)
		VALUES (?, ?, ?, ?, ?)
	`

	w := internal.Warehouse{
		ID:                 1,
		WarehouseCode:      "123ABC",
		Address:            "address",
		Telephone:          "telephone",
		MinimumCapacity:    1,
		MinimumTemperature: 1,
	}

	t.Run("case 1: success - Warehouse saved", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(w.WarehouseCode, w.Address, w.Telephone, w.MinimumCapacity, w.MinimumTemperature).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rp := repository.NewWarehouseMysqlRepository(db)
		err := rp.Save(&w)

		require.NoError(t, err)
	})

	t.Run("case 2: error - Error executing the query", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(w.WarehouseCode, w.Address, w.Telephone, w.MinimumCapacity, w.MinimumTemperature).
			WillReturnError(sql.ErrConnDone)

		rp := repository.NewWarehouseMysqlRepository(db)
		err := rp.Save(&w)

		require.Error(t, err)
	})

	t.Run("case 3: error - Error retrieving the last inserted ID", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(w.WarehouseCode, w.Address, w.Telephone, w.MinimumCapacity, w.MinimumTemperature).
			WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

		rp := repository.NewWarehouseMysqlRepository(db)
		err := rp.Save(&w)

		require.Error(t, err)
	})

}

func TestWarehouseMysql_Update(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	query := `
		UPDATE warehouses
		SET
			warehouse_code = ?, address = ?, telephone = ?, minimum_capacity = ?, minimum_temperature = ?
		WHERE
			id = ?;
	`

	w := internal.Warehouse{
		ID:                 1,
		WarehouseCode:      "123ABC",
		Address:            "address",
		Telephone:          "telephone",
		MinimumCapacity:    1,
		MinimumTemperature: 1,
	}

	t.Run("case 1: success - Warehouse updated", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(w.WarehouseCode, w.Address, w.Telephone, w.MinimumCapacity, w.MinimumTemperature, w.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rp := repository.NewWarehouseMysqlRepository(db)
		err := rp.Update(&w)

		require.NoError(t, err)
	})

	t.Run("case 2: error - Error executing the query", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(w.WarehouseCode, w.Address, w.Telephone, w.MinimumCapacity, w.MinimumTemperature, w.ID).
			WillReturnError(sql.ErrConnDone)

		rp := repository.NewWarehouseMysqlRepository(db)
		err := rp.Update(&w)

		require.Error(t, err)
	})
}

func TestWarehouseMysql_Delete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	query := `
		DELETE FROM warehouses
		WHERE
			id = ?;
	`

	t.Run("case 1: success - Warehouse deleted", func(t *testing.T) {
		id := 1

		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rp := repository.NewWarehouseMysqlRepository(db)
		err := rp.Delete(id)

		require.NoError(t, err)
	})

	t.Run("case 2: error - Error executing the query", func(t *testing.T) {
		id := 1

		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnError(sql.ErrConnDone)

		rp := repository.NewWarehouseMysqlRepository(db)
		err := rp.Delete(id)

		require.Error(t, err)
	})
}
