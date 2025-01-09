package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	ProductBatchAlreadyExists = errors.New("product-batch already exists")
)

func NewRepositoryProductBatchDB(db *sql.DB) *ProductBatchDB {
	return &ProductBatchDB{db}
}

type ProductBatchDB struct {
	db *sql.DB
}

func (r *ProductBatchDB) Save(prodBatch *internal.ProductBatch) error {
	result, err := r.db.Exec(
		"INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minumum_temperature, product_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		prodBatch.BatchNumber,
		prodBatch.CurrentQuantity,
		prodBatch.CurrentTemperature,
		prodBatch.DueDate,
		prodBatch.InitialQuantity,
		prodBatch.ManufacturingDate,
		prodBatch.ManufacturingHour,
		prodBatch.MinumumTemperature,
		prodBatch.ProductId,
		prodBatch.SectionId,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return ProductBatchAlreadyExists
			}
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	prodBatch.ID = int(id)

	return nil
}

func (r *ProductBatchDB) ProductBatchNumberExists(batchNumber int) (bool, error) {
	query := "SELECT COUNT(*) FROM product_batches WHERE batch_number = ?"

	var count int
	err := r.db.QueryRow(query, batchNumber).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
