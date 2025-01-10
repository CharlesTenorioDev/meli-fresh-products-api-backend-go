package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	ProductBatchNotFound            = errors.New("product-batch not found")
	ProductBatchAlreadyExists       = errors.New("product-batch already exists")
	ProductBatchNumberAlreadyInUse  = errors.New("product-batch with given product-batch number already registered")
	ProductBatchUnprocessableEntity = errors.New("couldn't parse product-batch")
)

func NewServiceProductBatch(rpProductBatch internal.ProductBatchRepository, rpSection internal.SectionRepository, rpProduct internal.ProductRepository) *ProductBatchService {
	return &ProductBatchService{
		rpB: rpProductBatch,
		rpS: rpSection,
		rpP: rpProduct,
	}
}

type ProductBatchService struct {
	rpB internal.ProductBatchRepository
	rpS internal.SectionRepository
	rpP internal.ProductRepository
}

func (s *ProductBatchService) FindByID(id int) (internal.ProductBatch, error) {
	prodBatch, err := s.rpB.FindByID(id)
	if err != nil {
		return internal.ProductBatch{}, ProductBatchNotFound
	}

	return prodBatch, nil
}

func (s *ProductBatchService) Save(prodBatch *internal.ProductBatch) error {
	if ok := prodBatch.Ok(); !ok {
		return ProductBatchUnprocessableEntity
	}

	countExists, err := s.rpB.ProductBatchNumberExists(prodBatch.BatchNumber)
	if err != nil || countExists {
		return ProductBatchNumberAlreadyInUse
	}

	_, err = s.rpS.FindByID(prodBatch.SectionId)
	if err != nil {
		return SectionNotFound
	}

	_, err = s.rpP.FindByID(prodBatch.ProductId)
	if err != nil {
		return ProductNotFound
	}

	err = s.rpB.Save(prodBatch)
	if err != nil {
		return err
	}

	return nil
}
