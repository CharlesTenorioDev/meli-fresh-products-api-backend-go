package internal

import "errors"

var (
	ProductTypeAlreadyExists = errors.New("product-type already exists")
	ProductTypeNotFound      = errors.New("product-type not found")
)

type ProductType struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type ProductTypeRepository interface {
	FindByID(id int) (ProductType, error)
}
