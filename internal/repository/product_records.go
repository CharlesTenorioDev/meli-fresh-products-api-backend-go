package repository

import (
	"database/sql"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
)

type ProductRecordsSQL struct {
	db *sql.DB
}

func NewProductRecordsSQL(db *sql.DB) *ProductRecordsSQL {
	return &ProductRecordsSQL{db}
}

func (psql *ProductRecordsSQL) FindAll() (productRecords []internal.ProductRecords, err error) {
	rows, err := psql.db.Query("SELECT `id`, `last_update_date`, `purchase_price`, `sale_price`, `product_id` FROM `product_records`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productRecord internal.ProductRecords

		err := rows.Scan(&productRecord.ID, &productRecord.LastUpdateDate, &productRecord.PurchasePrice, &productRecord.SalePrice, &productRecord.ProductID)
		if err != nil {
			return nil, err
		}

		productRecords = append(productRecords, productRecord)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return productRecords, nil
}

func (psql *ProductRecordsSQL) FindByID(id int) (internal.ProductRecords, error) {
	var productRecord internal.ProductRecords

	row := psql.db.QueryRow("SELECT `id`, `last_update_date`, `purchase_price`, `sale_price`, `product_id` FROM `product_records` WHERE `id` = ?", id)
	err := row.Scan(&productRecord.ID, &productRecord.LastUpdateDate, &productRecord.PurchasePrice, &productRecord.SalePrice, &productRecord.ProductID)

	if err != nil {
		return productRecord, service.ErrProductRecordsNotFound
	}

	return productRecord, nil
}

func (psql *ProductRecordsSQL) Save(productRec internal.ProductRecords) (internal.ProductRecords, error) {
	_, err := psql.db.Exec(
		"INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES (?, ?, ?, ?)",
		productRec.LastUpdateDate, productRec.PurchasePrice, productRec.SalePrice, productRec.ProductID,
	)

	if err != nil {
		return productRec, err
	}

	return productRec, nil
}
