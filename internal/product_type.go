package internal

import "errors"

var (
	ErrProductTypeAlreadyExists = errors.New("product-type already exists")
	ErrProductTypeNotFound      = errors.New("product-type not found")
)

type ProductType struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type ProductTypeRepository interface {
	FindByID(id int) (ProductType, error)
}
