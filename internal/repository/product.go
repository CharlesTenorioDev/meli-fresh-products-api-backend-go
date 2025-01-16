package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type ProductSQL struct {
	db *sql.DB
}

func NewProductSQL(db *sql.DB) *ProductSQL {
	return &ProductSQL{db}
}

const (
	FindAllString  = "SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products"
	FindByIdString = "SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products WHERE id = ?"
	SaveString     = "INSERT INTO products (id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	UpdateString   = `UPDATE product 
		 SET description = ?, expiration_rate = ?, freezing_rate = ?, 
		     height = ?, length = ?, net_weight = ?, 
		     product_code = ?, recommended_freezing_temperature = ?, 
		     width = ?, product_type_id = ?, seller_id = ?
		 WHERE id = ?`
	DeleteString         = "DELETE FROM products WHERE id = ?"
	FindAllRecordString  = "SELECT pr.product_id, p.description, COUNT(*) AS records_count FROM product_records pr JOIN products p ON pr.product_id = p.id GROUP BY pr.product_id, p.description;"
	FindByIdRecordString = "SELECT pr.product_id, p.description, COUNT(*) AS records_count FROM product_records pr JOIN products p ON pr.product_id = p.id WHERE p.id = ? GROUP BY pr.product_id, p.description;"
)

func (psql *ProductSQL) FindAll() (products []internal.Product, err error) {
	rows, err := psql.db.Query(FindAllString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrProductNotFound
		}
		return
	}
	defer rows.Close()
	for rows.Next() {
		var product internal.Product

		err := rows.Scan(&product.Id, &product.Description, &product.ExpirationRate, &product.FreezingRate,
			&product.Height, &product.Length, &product.NetWeight, &product.ProductCode, &product.RecommendedFreezingTemperature,
			&product.Width, &product.ProductTypeId, &product.SellerId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = internal.ErrProductNotFound
			}
		}
		products = append(products, product)
	}
	return
}

func (psql *ProductSQL) FindByID(id int) (internal.Product, error) {
	var product internal.Product

	row := psql.db.QueryRow(FindByIdString, id)
	err := row.Scan(&product.Id, &product.Description, &product.ExpirationRate, &product.FreezingRate,
		&product.Height, &product.Length, &product.NetWeight, &product.ProductCode, &product.RecommendedFreezingTemperature,
		&product.Width, &product.ProductTypeId, &product.SellerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrProductNotFound
		}
		return product, err
	}

	return product, nil
}

func (psql *ProductSQL) Save(product internal.Product) (p internal.Product, err error) {
	_, err = psql.db.Exec(
		SaveString,
		product.Id,
		product.Description,
		product.ExpirationRate,
		product.FreezingRate,
		product.Height,
		product.Length,
		product.NetWeight,
		product.ProductCode,
		product.RecommendedFreezingTemperature,
		product.Width,
		product.ProductTypeId,
		product.SellerId,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErroProductConflit
			}
		}
	}

	return
}

func (psql *ProductSQL) Update(product internal.Product) (internal.Product, error) {
	_, err := psql.db.Exec(
		UpdateString,
		product.Description,
		product.ExpirationRate,
		product.FreezingRate,
		product.Height,
		product.Length,
		product.NetWeight,
		product.ProductCode,
		product.RecommendedFreezingTemperature,
		product.Width,
		product.ProductTypeId,
		product.SellerId,
		product.Id,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErroProductConflit
			default:
				err = internal.ErrProductNotFound
			}
		}
		return product, err
	}

	return product, nil
}

func (psql *ProductSQL) Delete(id int) error {

	_, err := psql.db.Exec(DeleteString, id)
	if err != nil {
		return internal.ErrProductNotFound
	}

	return nil
}

func (psql *ProductSQL) FindAllRecord() ([]internal.ProductRecordsJsonCount, error) {
	var products []internal.ProductRecordsJsonCount

	rows, err := psql.db.Query(FindAllRecordString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var product internal.ProductRecordsJsonCount
		err := rows.Scan(&product.ProductID, &product.Description, &product.RecordsCount)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = internal.ErrProductNotFound
			}
			return products, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (psql *ProductSQL) FindByIdRecord(id int) (internal.ProductRecordsJsonCount, error) {
	var product internal.ProductRecordsJsonCount

	row := psql.db.QueryRow(FindByIdRecordString, id)
	err := row.Scan(&product.ProductID, &product.Description, &product.RecordsCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrProductIdNotFound
		}
		return product, err
	}

	return product, nil
}
