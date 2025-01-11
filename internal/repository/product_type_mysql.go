package repository

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewProductTypeMysql(db *sql.DB) *ProductTypeMysql {
	return &ProductTypeMysql{db}
}

type ProductTypeMysql struct {
	db *sql.DB
}

func (r *ProductTypeMysql) FindByID(id int) (internal.ProductType, error) {
	query := `
	SELECT 
		id,
		description          
	FROM 
		product_type
	WHERE 
		id = ?`

	var pt internal.ProductType
	err := r.db.QueryRow(query, id).Scan(
		&pt.ID,
		&pt.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pt, internal.ProductTypeNotFound
		}
		return pt, err
	}

	return pt, nil
}
