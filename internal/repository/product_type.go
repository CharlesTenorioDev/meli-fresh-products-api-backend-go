package repository

import (
	"errors"
	"fmt"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/loader"
)

func NewProductType() *ProductTypeDB {
	bdProdTypes := make(map[int]*internal.ProductType)

	prodTypeDBList, err := loader.ReadAllProductsTypeToFile()
	if err != nil {
		fmt.Println(err.Error())
		return &ProductTypeDB{
			DB: bdProdTypes,
		}
	}

	for _, value := range prodTypeDBList {
		section := internal.ProductType{
			ID: value.ID,
		}

		bdProdTypes[value.ID] = &section
	}

	return &ProductTypeDB{
		DB: bdProdTypes,
	}
}

type ProductTypeDB struct {
	DB map[int]*internal.ProductType
}

func (r *ProductTypeDB) FindByID(id int) (internal.ProductType, error) {
	product_type, exists := r.DB[id]
	if !exists {
		return internal.ProductType{}, errors.New("product_type not found")
	}
	return *product_type, nil
}
