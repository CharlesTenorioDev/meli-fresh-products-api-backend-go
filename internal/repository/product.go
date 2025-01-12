package repository

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type ProductSQL struct {
	db *sql.DB
}

func NewProductSQL(db *sql.DB) *ProductSQL {
	return &ProductSQL{db}
}

func (psql *ProductSQL) FindAll() (products []internal.Product, err error) {
	rows, err := psql.db.Query("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var product internal.Product
		err := rows.Scan(&product.Id, &product.Description, &product.ExpirationRate, &product.FreezingRate, &product.Height, &product.Length, &product.NetWeight, &product.ProductCode, &product.RecommendedFreezingTemperature, &product.Width, &product.ProductTypeId, &product.SellerId)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
func (psql *ProductSQL) FindByID(id int) (internal.Product, error) {
	var product internal.Product

	row := psql.db.QueryRow("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products WHERE id = ?", id)
	err := row.Scan(&product.Id, &product.Description, &product.ExpirationRate, &product.FreezingRate, &product.Height, &product.Length, &product.NetWeight, &product.ProductCode, &product.RecommendedFreezingTemperature, &product.Width, &product.ProductTypeId, &product.SellerId)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (psql *ProductSQL) Save(product internal.Product) (internal.Product, error) {
	_, err := psql.db.Exec(
		"INSERT INTO products (id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
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
		return product, err
	}

	return product, nil
}

func (psql *ProductSQL) Update(product internal.Product) (internal.Product, error) {
	result, err := psql.db.Exec(
		`UPDATE product 
		 SET description = ?, expiration_rate = ?, freezing_rate = ?, 
		     height = ?, length = ?, net_weight = ?, 
		     product_code = ?, recommended_freezing_temperature = ?, 
		     width = ?, product_type_id = ?, seller_id = ?
		 WHERE id = ?`,
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
		return internal.Product{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return internal.Product{}, err
	}

	if rowsAffected == 0 {
		return internal.Product{}, errors.New("product not found")
	}
	return product, nil
}

func (psql *ProductSQL) Delete(id int) error {
	result, err := psql.db.Exec("DELETE FROM product WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (psql *ProductSQL) FindAllRecord() ([]internal.ProductRecordsJsonCount, error) {
	var products []internal.ProductRecordsJsonCount

	rows, err := psql.db.Query("SELECT pr.product_id, p.description, COUNT(*) AS records_count FROM product_records pr JOIN products p ON pr.product_id = p.id GROUP BY pr.product_id, p.description;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var product internal.ProductRecordsJsonCount
		err := rows.Scan(&product.ProductID, &product.Description, &product.RecordsCount)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (psql *ProductSQL) FindByIdRecord(id int) (internal.ProductRecordsJsonCount, error) {
	var product internal.ProductRecordsJsonCount

	row := psql.db.QueryRow("SELECT pr.product_id, p.description, COUNT(*) AS records_count FROM product_records pr JOIN products p ON pr.product_id = p.id WHERE p.id = ? GROUP BY pr.product_id, p.description;", id)
	err := row.Scan(&product.ProductID, &product.Description, &product.RecordsCount)
	if err != nil {
		return product, err
	}

	return product, nil
}
