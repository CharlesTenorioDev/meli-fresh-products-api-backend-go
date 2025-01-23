package repository

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

const FindByIDProductType = `
SELECT 
	id,
	description          
FROM 
	product_type
WHERE 
	id = ?`

func NewProductTypeMysql(db *sql.DB) *ProductTypeMysql {
	return &ProductTypeMysql{db}
}

type ProductTypeMysql struct {
	db *sql.DB
}

func (r *ProductTypeMysql) FindByID(id int) (internal.ProductType, error) {
	var pt internal.ProductType
	err := r.db.QueryRow(FindByIDProductType, id).Scan(
		&pt.ID,
		&pt.Description,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pt, internal.ErrProductTypeNotFound
		}

		return pt, err
	}

	return pt, nil
}
