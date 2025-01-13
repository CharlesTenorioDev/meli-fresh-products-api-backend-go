package service

import (
	"fmt"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
)

type ProductRecordsDefault struct {
	productRecRepo internal.ProductRecordsRepository
	productRepo    internal.ProductRepository
}

func NewProductRecordsDefault(prodRecRepo internal.ProductRecordsRepository, prodRepo internal.ProductRepository) *ProductRecordsDefault {
	return &ProductRecordsDefault{
		productRecRepo: prodRecRepo,
		productRepo:    prodRepo,
	}
}

// Implementação da função Create
func (pr *ProductRecordsDefault) Create(productRec internal.ProductRecords) (internal.ProductRecords, error) {
	if err := ValidateProductRec(productRec); err != nil {
		return productRec, err
	}
	product, err := pr.productRepo.FindByID(productRec.ProductID)
	fmt.Printf("Product found: %+v\n", product)
	if err != nil {
		return productRec, ProductNotExists
	}
	return pr.productRecRepo.Save(productRec)
}

// Implementação da função GetAll
func (pr *ProductRecordsDefault) GetAll() ([]internal.ProductRecords, error) {
	productRecords, err := pr.productRecRepo.FindAll()
	if err != nil {
		return nil, rest_err.NewInternalServerError("Erro ao buscar todos os registros de produtos")
	}
	return productRecords, nil
}

// Implementação da função GetByID
func (pr *ProductRecordsDefault) GetByID(id int) (internal.ProductRecords, error) {
	if id <= 0 {
		return internal.ProductRecords{}, rest_err.NewBadRequestError("O ID deve ser válido e maior que zero")
	}

	productRecord, err := pr.productRecRepo.FindByID(id)
	if err != nil {
		return internal.ProductRecords{}, err
	}

	return productRecord, nil
}

func ValidateProductRec(productRec internal.ProductRecords) error {
	if productRec.LastUpdateDate.IsZero() ||
		productRec.PurchasePrice <= 0 ||
		productRec.SalePrice <= 0 {
		return ProductUnprocessableEntity
	}
	return nil
}
