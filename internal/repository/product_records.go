package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type ProductRecordsSQL struct {
	db *sql.DB
}

func NewProductRecordsSQL(db *sql.DB) *ProductRecordsSQL {
	return &ProductRecordsSQL{db}
}

const (
	FindAllProductRecords  = "SELECT `id`, `last_update_date`, `purchase_price`, `sale_price`, `product_id` FROM `product_records`"
	FindByIdProductRecords = "SELECT `id`, `last_update_date`, `purchase_price`, `sale_price`, `product_id` FROM `product_records` WHERE `id` = ?"
	SaveProductRecords     = "INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES (?, ?, ?, ?)"
)

// Implementação de FindAll
func (psql *ProductRecordsSQL) FindAll() (productRecords []internal.ProductRecords, err error) {
	rows, err := psql.db.Query(FindAllProductRecords)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productRecord internal.ProductRecords

		err := rows.Scan(&productRecord.ID, &productRecord.LastUpdateDate, &productRecord.PurchasePrice, &productRecord.SalePrice, &productRecord.ProductID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = internal.ErrProductNotFound
			}
			return productRecords, err
		}

		productRecords = append(productRecords, productRecord)
	}

	return
}

func (psql *ProductRecordsSQL) FindByID(id int) (internal.ProductRecords, error) {
	var productRecord internal.ProductRecords

	row := psql.db.QueryRow(FindByIdProductRecords, id)
	err := row.Scan(&productRecord.ID, &productRecord.LastUpdateDate, &productRecord.PurchasePrice, &productRecord.SalePrice, &productRecord.ProductID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrProductReordsNotFound
		}
		return productRecord, err
	}

	return productRecord, nil
}

func (psql *ProductRecordsSQL) Save(productRec internal.ProductRecords) (internal.ProductRecords, error) {
	_, err := psql.db.Exec(
		SaveProductRecords,
		productRec.LastUpdateDate, productRec.PurchasePrice, productRec.SalePrice, productRec.ProductID,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErrProductReordsConflict
			}
		}
		return productRec, err
	}

	return productRec, nil
}
