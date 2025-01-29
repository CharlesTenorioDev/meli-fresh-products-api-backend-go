package repository

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewWarehouseMysqlRepository(db *sql.DB) *WarehouseMysqlRepository {
	return &WarehouseMysqlRepository{db}
}

type WarehouseMysqlRepository struct {
	db *sql.DB
}

func (w *WarehouseMysqlRepository) FindAll() ([]internal.Warehouse, error) {
	warehouses := make([]internal.Warehouse, 0)

	query := `
		SELECT
			id, warehouse_code, address, telephone, minimum_capacity, minimum_temperature
		FROM
			warehouses;
	`

	/// executing the query
	rows, err := w.db.Query(query)
	if err != nil {
		return warehouses, err
	}
	defer rows.Close()

	// iterating over the rows
	for rows.Next() {
		var warehouse internal.Warehouse
		err := rows.Scan(
			&warehouse.ID,
			&warehouse.WarehouseCode,
			&warehouse.Address,
			&warehouse.Telephone,
			&warehouse.MinimumCapacity,
			&warehouse.MinimumTemperature,
		)
		if err != nil {
			return nil, err
		}
		// appending the warehouse to the slice
		warehouses = append(warehouses, warehouse)
	}
	err = rows.Err()
	if err != nil {
		return warehouses, err
	}

	return warehouses, err
}

func (w *WarehouseMysqlRepository) FindByID(id int) (internal.Warehouse, error) {
	query := `
		SELECT
			id, warehouse_code, address, telephone, minimum_capacity, minimum_temperature
		FROM
			warehouses
		WHERE
			id = ?;
	`
	// creating a new warehouse
	var warehouse internal.Warehouse
	// executing the query
	err := w.db.QueryRow(query, id).Scan(
		&warehouse.ID,
		&warehouse.WarehouseCode,
		&warehouse.Address,
		&warehouse.Telephone,
		&warehouse.MinimumCapacity,
		&warehouse.MinimumTemperature,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return warehouse, internal.ErrWarehouseRepositoryNotFound
		}

		return warehouse, err
	}

	return warehouse, nil
}

func (w *WarehouseMysqlRepository) Save(warehouse *internal.Warehouse) error {
	query := `
		INSERT INTO warehouses (warehouse_code, address, telephone, minimum_capacity, minimum_temperature)
		VALUES (?, ?, ?, ?, ?)
	`
	// executing the query
	result, err := w.db.Exec(
		query,
		(*warehouse).WarehouseCode,
		(*warehouse).Address,
		(*warehouse).Telephone,
		(*warehouse).MinimumCapacity,
		(*warehouse).MinimumTemperature,
	)
	if err != nil {
		return err
	}
	// getting the ID of the last inserted warehouse
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// setting the ID of the warehouse
	(*warehouse).ID = int(id)

	return nil
}

func (w *WarehouseMysqlRepository) Update(warehouse *internal.Warehouse) error {
	query := `
		UPDATE warehouses
		SET
			warehouse_code = ?, address = ?, telephone = ?, minimum_capacity = ?, minimum_temperature = ?
		WHERE
			id = ?;
	`

	// executing the query
	_, err := w.db.Exec(
		query,
		warehouse.WarehouseCode,
		warehouse.Address,
		warehouse.Telephone,
		warehouse.MinimumCapacity,
		warehouse.MinimumTemperature,
		warehouse.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (w *WarehouseMysqlRepository) Delete(id int) error {
	query := `
		DELETE FROM warehouses
		WHERE
			id = ?;
	`

	// executing the query
	_, err := w.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
