package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

type ProductRecordsDefault struct {
	productRecRepo internal.ProductRecordsRepository
	productRepo    internal.ProductRepository
}

var (
	ErrProductRecordsNotFound = errors.New("product records not found")
)

func NewProductRecordsDefault(prodRecRepo internal.ProductRecordsRepository, prodRepo internal.ProductRepository) *ProductRecordsDefault {
	return &ProductRecordsDefault{
		productRecRepo: prodRecRepo,
		productRepo:    prodRepo,
	}
}

func (pr *ProductRecordsDefault) Create(productRec internal.ProductRecords) (internal.ProductRecords, error) {
	if err := ValidateProductRec(productRec); err != nil {
		return productRec, err
	}

	_, err := pr.productRepo.FindByID(productRec.ProductID)
	if err != nil {
		return productRec, internal.ErrProductIdNotFound
	}

	return pr.productRecRepo.Save(productRec)
}

func (pr *ProductRecordsDefault) GetAll() ([]internal.ProductRecords, error) {
	productRecords, err := pr.productRecRepo.FindAll()
	if err != nil {
		return nil, resterr.NewInternalServerError("Erro ao buscar todos os registros de produtos")
	}

	return productRecords, nil
}

func (pr *ProductRecordsDefault) GetByID(id int) (internal.ProductRecords, error) {
	if id <= 0 {
		return internal.ProductRecords{}, resterr.NewBadRequestError("O ID deve ser válido e maior que zero")
	}

	productRecord, err := pr.productRecRepo.FindByID(id)
	if err != nil {
		return internal.ProductRecords{}, err
	}

	return productRecord, nil
}

func ValidateProductRec(productRec internal.ProductRecords) error {
	if productRec.LastUpdateDate.IsZero() {
		return errors.New("LastUpdateDate is empty.")
	}
	if productRec.PurchasePrice <= 0 {
		return errors.New("PurchasePrice is empty.")
	}
	if productRec.SalePrice <= 0 {
		return errors.New("SalePrice is empty.")
	}

	return nil
}

